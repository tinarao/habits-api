package habits

import (
	"fmt"
	"hbapi/internal/db"

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
