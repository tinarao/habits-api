package main

import (
	"hbapi/internal/db"
	"hbapi/modules/auth"
	"hbapi/router"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to preload .env vars", "fatal", err.Error())
		os.Exit(1)
	}

	srv := gin.Default()

	db.Init()
	auth.Init()
	router.SetupRoutes(srv)

	port := os.Getenv("PORT")
	err := srv.Run(port)
	if err != nil {
		slog.Error("failed to start a server", "fatal", err.Error())
		os.Exit(1)
	}
}
