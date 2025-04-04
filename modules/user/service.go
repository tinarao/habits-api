package user

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"

	"github.com/markbates/goth"
)

func CompleteAuth(u *goth.User) (profile *db.User, err error) {
	account, exists := FindByEmail(u.Email)
	if !exists {
		acc, err := PersistAccount(u)
		if err != nil {
			return nil, err
		}

		return acc, nil
	}

	return account, nil
}

func FindByEmail(email string) (profile *db.User, exists bool) {
	user := &db.User{}
	dbr := db.Client.Where("email = ?", email).First(&user)
	if dbr.Error != nil {
		slog.Error("db error", "error", dbr.Error.Error())
		return nil, false
	}

	return user, true
}

func PersistAccount(u *goth.User) (profile *db.User, err error) {
	if u == nil || u.Email == "" {
		return nil, fmt.Errorf("invalid user")
	}

	user := &db.User{
		Plan:                db.FreePlan,
		Provider:            u.Provider,
		Name:                &u.Name,
		Email:               u.Email,
		ImageUrl:            &u.AvatarURL,
		RefreshToken:        &u.RefreshToken,
		RefreshTokenExpires: u.ExpiresAt,
	}

	dbr := db.Client.Create(user)
	if dbr.Error != nil {
		return nil, err
	}

	return user, nil
}
