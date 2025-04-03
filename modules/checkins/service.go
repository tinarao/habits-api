package checkins

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"
	"time"
)

func Create(u *db.User, habitId uint) (*db.CheckIn, error) {
	if u == nil || habitId == 0 {
		return nil, fmt.Errorf("invalid user or habit id")
	}

	s := db.CheckIn{
		UserId:  u.ID,
		HabitId: habitId,
	}

	dbr := db.Client.Create(&s)
	if dbr.Error != nil {
		slog.Error("failed to create checkin", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to create checkin")
	}

	return &s, nil
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

// Returns checkins made in previous 2 days
func GetLatest() ([]*db.CheckIn, error) {
	checkins := make([]*db.CheckIn, 0)
	now := time.Now()
	twoDaysEarlier := now.Add(time.Hour * -48)
	dbr := db.Client.Where("created_at >= ? AND created_at <= ?", twoDaysEarlier, now).Find(checkins)
	if dbr.Error != nil {
		slog.Error("failed to find latest checkins", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to find latest checkins")
	}

	return checkins, nil
}
