package router

import (
	"hbapi/modules/auth"
	"hbapi/modules/habits"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/hc", func(c *gin.Context) {
		c.String(http.StatusOK, "Alive")
	})

	// Modules
	api := r.Group("/api")
	auth.SetupRoutes(api)
	habits.SetupRoutes(api)
}
