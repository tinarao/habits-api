package utils

import (
	"fmt"
	"hbapi/internal/db"
	"hbapi/modules/auth"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func SetupTestEnv() {
	rootpath, err := GetRootPath()
	if err != nil {
		panic(fmt.Sprintf("failed to get root path: %v", err))
	}

	envFilepath := filepath.Join(rootpath, ".env")
	if err := godotenv.Load(envFilepath); err != nil {
		panic(fmt.Sprintf("failed to load env file: %v", err))
	}

	db.Init()
	auth.Init()
}

func GetRootPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	root := filepath.Dir(filepath.Dir(filename))
	return root, nil
}
