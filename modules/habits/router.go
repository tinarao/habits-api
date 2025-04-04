package habits

import (
	"hbapi/modules/auth"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	habits := r.Group("/habits")

	habits.GET("/", auth.SessionsMiddleware, getAll)
	habits.GET("/slug/:slug", auth.SessionsMiddleware, getBySlug)
	habits.GET("/pin", auth.SessionsMiddleware, getPinned)
	habits.POST("/", auth.SessionsMiddleware, create)
	habits.PATCH("/pin/:slug", auth.SessionsMiddleware, togglePin)
	habits.PATCH("/rename/:slug/:newName", auth.SessionsMiddleware, rename)
	habits.DELETE("/:slug", auth.SessionsMiddleware, delete)
}

//

func create(c *gin.Context) {
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dto := &createDTO{}

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect body data shape"})
		return
	}

	created, err := Create(*dto, user)
	if err != nil {
		slog.Error("failed to create habit", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "habit": *created})
}

func getAll(c *gin.Context) {
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	habits, err := GetAll(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Failed to retrieve habits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habits": habits})
}

func getBySlug(c *gin.Context) {
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	slug := c.Param("slug")
	habit, err := GetBySlug(slug, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Failed to find specified habit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habit": habit})
}

func getPinned(c *gin.Context) {
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	habits, err := GetPinned(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "failed to get pinned habits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habits": habits})
}

func delete(c *gin.Context) {
	slug := c.Param("slug")
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = Delete(slug, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func rename(c *gin.Context) {
	slug := c.Param("slug")
	newName := c.Param("newName")
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = Rename(slug, user, newName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func togglePin(c *gin.Context) {
	slug := c.Param("slug")
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = TogglePin(slug, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusOK)
}
