package paywall

import (
	"errors"
	"fmt"
	"hbapi/internal/db"
)

func ProtectCreate(u *db.User) error {
	if u.Plan == db.FreePlan {
		var count int64 = 0
		_ = db.Client.Model(&db.Habit{}).Where("user_id = ?", u.ID).Count(&count)
		fmt.Printf("count: %d\n", count)

		if count >= 5 {
			return errors.New("not allowed")
		}
	}

	return nil
}

func ProtectRemind(u *db.User) error {
	if u.Plan == db.FreePlan {
		return errors.New("forbidden due to subscription plan")
	}

	return nil
}
