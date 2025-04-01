package auth

import (
	"hbapi/modules/user"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func SetupRoutes(r *gin.RouterGroup) {
	oauth := r.Group("/auth")

	oauth.GET("/logout", logout)
	oauth.GET("/login/:provider", login)
	oauth.GET("/callback/:provider", callback)
}

func login(c *gin.Context) {
	provider := c.Param("provider")
	if !slices.Contains(AvailableProviders, provider) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Provider is unknown or temporary unavailable",
		})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func logout(c *gin.Context) {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		slog.Error("failed to logout", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to destroy session"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func callback(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	u, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		slog.Error("failed authorization attempt", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}

	_, err = user.CompleteAuth(&u)
	if err != nil {
		slog.Error("failed to complete authorization", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}

	session, _ := store.Get(c.Request, "session")
	session.Values["user"] = u
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		slog.Error("failed to save the session", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/dashboard")
}
