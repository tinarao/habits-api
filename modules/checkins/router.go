package checkins

import (
	"hbapi/modules/auth"
	"hbapi/modules/habits"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	checkins := r.Group("/checkins")

	checkins.GET("/habit/:habitSlug", auth.SessionsMiddleware, getByHabitSlug)
	checkins.GET("/latest", auth.SessionsMiddleware, getLatest)
	checkins.POST("/:habitSlug", auth.SessionsMiddleware, create)
}

func create(c *gin.Context) {
	// insert user shit here
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	habitSlug := c.Param("habitSlug")
	habit, err := habits.GetBySlug(habitSlug, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Habit is not found"})
		return
	}

	if habit.UserId != user.ID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Habit is not found"})
	}

	_, err = Create(user, habit.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created successfully"})
}

func getByHabitSlug(c *gin.Context) {
	user, err := auth.GetUserFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	habitSlug := c.Param("habitSlug")
	habit, err := habits.GetBySlug(habitSlug, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Habit is not found"})
		return
	}

	if habit.UserId != user.ID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Habit is not found"})
		return
	}

	checkins, err := GetByHabit(habit.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Habit and/or checkins could not be found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkins": checkins})
}

// func getByDate() {}
func getLatest(c *gin.Context) {
	checkins, err := GetLatest()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Latest checkins can not be found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkins": checkins})
}
