package auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yandex"
	"github.com/michaeljs1990/sqlitestore"
)

const (
	COOKIE_MAX_AGE = time.Hour * 24 * 30
	COOKIE_PATH    = "/"
)

var (
	AvailableProviders = []string{"yandex"}
	store              *sqlitestore.SqliteStore
)

func Init() {
	secret := os.Getenv("SECRET")
	port := os.Getenv("PORT")
	dbpath := os.Getenv("DB_PATH")
	if secret == "" || port == "" || dbpath == "" {
		slog.Error("SECRET or PORT .env var is not set")
		os.Exit(1)
	}

	oauthRedirectPath := fmt.Sprintf("http://localhost%s/api/auth/callback", os.Getenv("PORT"))

	var err error
	store, err = sqlitestore.NewSqliteStore(dbpath, "sessions", COOKIE_PATH, int(COOKIE_MAX_AGE.Seconds()))
	if err != nil {
		slog.Error("failed to create a sessions store", "fatal", err.Error())
		os.Exit(1)
	}

	gothic.Store = store

	ydxSecret := os.Getenv("YANDEX_SECRET")
	ydxClientKey := os.Getenv("YANDEX_CLIENT_KEY")
	if ydxClientKey == "" || ydxSecret == "" {
		slog.Error("yandex outh-related vars are not set in .env")
		os.Exit(1)
	}

	goth.UseProviders(
		yandex.New(ydxClientKey, ydxSecret, fmt.Sprintf("%s/yandex", oauthRedirectPath)),
	)
}
