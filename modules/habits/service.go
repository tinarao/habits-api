package habits

import (
	"fmt"
	"hbapi/internal/db"
	"log/slog"

	"github.com/gosimple/slug"
)

func Create(dto createDTO, u *db.User) (*db.Habit, error) {
	slug := slug.Make(dto.Name)
	if dto.Name == "" {
		return nil, fmt.Errorf("invalid habit name")
	}
	habit := &db.Habit{
		Name:        dto.Name,
		Description: dto.Description,
		IsPinned:    false,
		Slug:        slug,
		User:        *u,
		UserId:      u.ID,
	}

	dbr := db.Client.Create(habit)
	if dbr.Error != nil {
		return nil, fmt.Errorf("failed to create a habit")
	}

	return habit, nil
}

func GetAll(u *db.User) ([]db.Habit, error) {
	var habits []db.Habit
	dbr := db.Client.Preload("CheckIns").Where("user_id = ?", u.ID).Find(&habits)
	if dbr.Error != nil {
		return nil, fmt.Errorf("failed to fetch habits")
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

func GetPinned(user *db.User) ([]db.Habit, error) {
	var habits []db.Habit
	dbr := db.Client.Where("is_pinned = ?", true).Find(&habits)
	if dbr.Error != nil {
		return nil, fmt.Errorf("failed to get pinned habits")
	}

	return habits, nil
}

func Rename(hSlug string, user *db.User, newName string) error {
	habit, err := GetBySlug(hSlug, user)
	if err != nil {
		return err
	}

	s := slug.Make(newName)
	habit.Name = newName
	habit.Slug = s

	dbr := db.Client.Save(&habit)
	if dbr.Error != nil {
		slog.Error("failed to rename a habit", "error", dbr.Error.Error())
		return fmt.Errorf("failed to remove a habit")
	}

	return nil
}

func TogglePin(slug string, u *db.User) error {
	habit, err := GetBySlug(slug, u)
	if err != nil {
		return err
	}

	habit.IsPinned = !habit.IsPinned
	dbr := db.Client.Save(&habit)
	if dbr.Error != nil {
		slog.Error("failed to toggle pin on habit", "error", dbr.Error.Error())
		return fmt.Errorf("failed to toggle pin")
	}

	return nil
}

func Delete(slug string, user *db.User) error {
	habit, err := GetBySlug(slug, user)
	if err != nil {
		return err
	}

	dbr := db.Client.Delete(&db.Habit{}, habit.ID)
	if dbr.Error != nil {
		slog.Error("failed to delete habit", "error", dbr.Error.Error())
		return fmt.Errorf("failed to delete habit")
	}

	return nil
}
