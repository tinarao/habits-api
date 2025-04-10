package habits

import (
	"errors"
	"hbapi/internal/db"
	"hbapi/internal/paywall"
	"log/slog"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func IsUserAllowedToCreate(u *db.User) bool {
	var count int64 = 0
	db.Client.Model(&db.Habit{}).Where("user_id = ?").Count(&count)
	if u.Plan == db.FreePlan && count >= 5 {
		return false
	}

	return true
}

func Create(dto createDTO, u *db.User) (*db.Habit, error) {
	err := paywall.ProtectCreate(u)
	if err != nil {
		return nil, err
	}

	uuid, _ := gonanoid.New()

	habit := &db.Habit{
		Name:        dto.Name,
		Description: dto.Description,
		IsPinned:    false,
		Slug:        uuid,
		User:        *u,
		UserId:      u.ID,
		Remind:      dto.Remind,
	}

	dbr := db.Client.Create(habit)
	if dbr.Error != nil {
		return nil, errors.New("failed to create a habit")
	}

	return habit, nil
}

func ChangeColor(habitSlug, newColor string, u *db.User) error {
	habit := &db.Habit{}
	dbr := db.Client.Where("slug = ? AND user_id = ?", habitSlug, u.ID).First(&habit)
	if dbr.Error != nil {
		return errors.New("habit does not exist")
	}

	habit.CheckinsColor = newColor
	dbr = db.Client.Save(&habit)
	if dbr.Error != nil {
		slog.Error("failed to change habit color", "error", dbr.Error.Error())
		return errors.New("failed to change habit color")
	}

	return nil
}

func GetRandomHabit(u *db.User) (*db.Habit, error) {
	habit := &db.Habit{}

	dbr := db.Client.
		Select("habits.*, COUNT(check_ins.id) as checkin_count").
		Joins("LEFT JOIN check_ins ON habits.id = check_ins.habit_id").
		Where("habits.user_id = ?", u.ID).
		Group("habits.id").
		Order("checkin_count DESC").
		Preload("CheckIns").
		First(&habit)
	if dbr.Error != nil {
		slog.Error("failed to retrieve random habit", "error", dbr.Error.Error())
		return nil, errors.New("failed to get random habit")
	}

	return habit, nil
}

func GetAll(u *db.User) ([]db.Habit, error) {
	var habits []db.Habit
	dbr := db.Client.Preload("CheckIns").Where("user_id = ?", u.ID).Find(&habits)
	if dbr.Error != nil {
		return nil, errors.New("failed to fetch habits")
	}

	return habits, nil
}

func GetBySlug(slug string, user *db.User) (*db.Habit, error) {
	habit := &db.Habit{}
	dbr := db.Client.Preload("CheckIns").Where("slug = ? AND user_id = ?", slug, user.ID).First(&habit)
	if dbr.Error != nil {
		return nil, errors.New("habit not found")
	}

	return habit, nil
}

func GetPinned(user *db.User) ([]db.Habit, error) {
	var habits []db.Habit
	dbr := db.Client.Where("is_pinned = ?", true).Find(&habits)
	if dbr.Error != nil {
		return nil, errors.New("failed to get pinned habits")
	}

	return habits, nil
}

func Rename(hSlug string, user *db.User, newName string) error {
	habit, err := GetBySlug(hSlug, user)
	if err != nil {
		return err
	}

	habit.Name = newName

	dbr := db.Client.Save(&habit)
	if dbr.Error != nil {
		slog.Error("failed to rename a habit", "error", dbr.Error.Error())
		return errors.New("failed to remove a habit")
	}

	return nil
}

func ToggleRemind(slug string, u *db.User) error {
	habit, err := GetBySlug(slug, u)
	if err != nil {
		return err
	}

	err = paywall.ProtectRemind(u)
	if err != nil {
		return err
	}

	habit.Remind = !habit.Remind
	dbr := db.Client.Save(&habit)
	if dbr.Error != nil {
		slog.Error("failed to toggle remind on habit", "error", dbr.Error.Error())
		return errors.New("failed to toggle remind")
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
		return errors.New("failed to toggle pin")
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
		return errors.New("failed to delete habit")
	}

	return nil
}

// Habits without checkin where date is today
func GetUncheckedHabits(localDate string, user *db.User) ([]db.Habit, error) {
	dateOnly := localDate[:10]

	var habits []db.Habit
	dbr := db.Client.
		Joins("LEFT JOIN check_ins ON habits.id = check_ins.habit_id AND DATE(check_ins.created_at) = ?", dateOnly).
		Where("check_ins.id IS NULL AND habits.user_id = ?", user.ID).
		Find(&habits)

	if dbr.Error != nil {
		slog.Error("failed to fetch unchecked habits", "error", dbr.Error.Error())
		return nil, errors.New("failed to fetch unchecked habits")
	}

	return habits, nil
}

// GetMostCheckedHabits returns habits sorted by check-in count in descending order
func GetMostCheckedHabits(user *db.User, limit int) ([]db.Habit, error) {
	var habits []db.Habit
	dbr := db.Client.
		Select("habits.*, COUNT(check_ins.id) as checkin_count").
		Joins("LEFT JOIN check_ins ON habits.id = check_ins.habit_id").
		Where("habits.user_id = ?", user.ID).
		Group("habits.id").
		Order("checkin_count DESC").
		Limit(limit).
		Preload("CheckIns").
		Find(&habits)

	if dbr.Error != nil {
		slog.Error("failed to fetch most checked habits", "error", dbr.Error.Error())
		return nil, errors.New("failed to fetch most checked habits")
	}

	return habits, nil
}
