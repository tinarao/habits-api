package main

import (
	"hbapi/internal/db"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to preload .env vars", "fatal", err.Error())
		os.Exit(1)
	}

	db.Init()

	srv := gin.Default()

	srv.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Yee")
	})

	port := os.Getenv("PORT")
	err := srv.Run(port)
	if err != nil {
		slog.Error("failed to start a server", "fatal", err.Error())
		os.Exit(1)
	}
}
