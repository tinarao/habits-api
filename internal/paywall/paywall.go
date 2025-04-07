package paywall

import (
	"errors"
	"hbapi/internal/db"
)

func ProtectCreate(u *db.User) error {
	if u.Plan == db.FreePlan {
		var count int64 = 0
		_ = db.Client.Where("user_id = ?", u.ID).Count(&count)

		if count >= 5 {
			return errors.New("not allowed")
		}
	}

	return nil
}
