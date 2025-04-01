package router

import (
	"hbapi/modules/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/hc", func(c *gin.Context) {
		c.String(http.StatusOK, "Alive")
	})

	api := r.Group("/api")
	auth.SetupRoutes(api)
}
