package habits

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"

	"github.com/gosimple/slug"
)

func Create(dto createDTO, u *db.User) (*db.Habit, error) {
	slug := slug.Make(dto.Name)

	habit := &db.Habit{
		Name:        dto.Name,
		Description: dto.Description,
		IsPinned:    false,
		UserId:      u.ID,
		Slug:        slug,
	}

	dbr := db.Client.Create(habit)
	if dbr.Error != nil {
		return nil, fmt.Errorf("failed to create a habit")
	}

	return habit, nil
}

func GetAll(u *db.User) ([]db.Habit, error) {
	habits := make([]db.Habit, 0)
	dbr := db.Client.Where("user_id = ?", u.ID).Find(&habits)
	if dbr.Error != nil {
		slog.Error("failed to retrieve []habit", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failedailed to retrieve \"%s\" habits", *u.NickName)
	}

	return habits, nil
}

func GetBySlug(slug string, user *db.User) (*db.Habit, error) {
	habit := &db.Habit{}
	dbr := db.Client.Where("slug = ? AND user_id = ?", slug, user.ID).First(&habit)
	if dbr.Error != nil {
		return nil, fmt.Errorf("habit not found")
	}

	return habit, nil
}
