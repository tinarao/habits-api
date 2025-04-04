package auth

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	port := os.Getenv("PORT")
	dbpath := os.Getenv("DB_PATH")
	secret := os.Getenv("SECRET")
	if port == "" || dbpath == "" || secret == "" {
		slog.Error("SECRET or PORT .env var is not set")
		os.Exit(1)
	}

	oauthRedirectPath := fmt.Sprintf("http://localhost%s/api/auth/callback", os.Getenv("PORT"))

	var err error
	store, err = sqlitestore.NewSqliteStore(
		dbpath,
		"sessions",
		COOKIE_PATH,
		int(COOKIE_MAX_AGE.Seconds()),
		[]byte(secret),
	)
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

func GetUserFromSession(c *gin.Context) (user *db.User, err error) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return nil, err
	}

	data, ok := session.Values["user"]
	if !ok {
		return nil, fmt.Errorf("failed to parse a session")
	}

	su := data.(goth.User)

	u := &db.User{}
	dbr := db.Client.Where("email = ?", su.Email).First(&u)
	if dbr.Error != nil {
		slog.Error("failed to find a user", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to find a user: %s", dbr.Error.Error())
	}

	return u, nil
}

func GetUserFromCtx(c *gin.Context) (user *db.User, err error) {
	u, ok := c.Get("user")
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}

	ctxUser := u.(*db.User)
	return ctxUser, nil
}
