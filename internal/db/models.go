package db

import "time"

type Plan string

const (
	FreePlan     Plan = "free"
	AdvancedPlan Plan = "advanced"
)

type Habit struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Remind      bool      `json:"remind" gorm:"default:false"`
	Description *string   `json:"description"`
	CheckIns    []CheckIn `json:"checkIns" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsPinned    bool      `json:"isPinned"`
	User        User      `json:"user"`
	UserId      uint      `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CheckIn struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	Habit     Habit
	HabitId   uint
	User      User
	UserId    uint
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	Habits              []Habit   `json:"habits"`
	CheckIns            []CheckIn `json:"checkIns"`
	Plan                Plan      `json:"plan" gorm:"default:free"`
	Provider            string    `json:"provider"`
	NickName            *string   `json:"nickname"`
	Name                *string   `json:"name"`
	Email               string    `json:"email"`
	ImageUrl            *string   `json:"imageUrl"`
	RefreshToken        *string   `json:"refreshToken"`
	RefreshTokenExpires time.Time `json:"refreshTokenExpires"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
