package habits

import (
	"hbapi/modules/auth"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	habits := r.Group("/habits")

	habits.POST("/", auth.SessionsMiddleware, create)
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
