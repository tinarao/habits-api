package checkins

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"
)

func Create(u *db.User, habitId uint) (*db.CheckIn, error) {
	s := &db.CheckIn{
		UserId:  u.ID,
		HabitId: habitId,
	}

	dbr := db.Client.Create(&s)
	if dbr.Error != nil {
		slog.Error("failed to create checkin", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to create checkin")
	}

	return s, nil
}

func GetByHabit(habitId uint) ([]*db.CheckIn, error) {
	checkins := make([]*db.CheckIn, 0)
	dbr := db.Client.Where("habit_id = ?", habitId).Find(&checkins)
	if dbr.Error != nil {
		slog.Error("failed to retrieve checkins", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to retrieve checkins")
	}

	return checkins, nil
}
