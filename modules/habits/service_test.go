package habits

import (
	"hbapi/internal/db"
	"hbapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testUser = &db.User{
	ID:       1,
	NickName: stringPtr("testUser"),
}

func stringPtr(s string) *string {
	return &s
}

func TestCreate(t *testing.T) {
	desc := "Test Description"
	tests := []struct {
		name    string
		dto     createDTO
		user    *db.User
		wantErr bool
	}{
		{
			name: "valid habit",
			dto: createDTO{
				Name:        "Test Habit",
				Description: &desc,
			},
			user:    testUser,
			wantErr: false,
		},
		{
			name: "empty name",
			dto: createDTO{
				Name:        "",
				Description: &desc,
			},
			user:    testUser,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			habit, err := Create(tt.dto, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, habit)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, habit)
				assert.Equal(t, tt.dto.Name, habit.Name)
				assert.Equal(t, tt.dto.Description, habit.Description)
				assert.Equal(t, tt.user.ID, habit.UserId)
				assert.False(t, habit.IsPinned)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name    string
		user    *db.User
		setup   func()
		wantErr bool
	}{
		{
			name: "user with habits",
			user: testUser,
			setup: func() {
				Create(createDTO{Name: "Habit 1"}, testUser)
				Create(createDTO{Name: "Habit 2"}, testUser)
			},
			wantErr: false,
		},
		{
			name:    "user without habits",
			user:    testUser,
			setup:   func() {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			habits, err := GetAll(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, habits)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, habits)
				for _, habit := range habits {
					assert.Equal(t, tt.user.ID, habit.UserId)
				}
			}
		})
	}
}

func TestGetBySlug(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		user    *db.User
		setup   func() *db.Habit
		wantErr bool
	}{
		{
			name: "existing habit",
			slug: "test-habit",
			user: testUser,
			setup: func() *db.Habit {
				habit, _ := Create(createDTO{Name: "Test Habit"}, testUser)
				return habit
			},
			wantErr: false,
		},
		{
			name:    "non-existing habit",
			slug:    "non-existing",
			user:    testUser,
			setup:   func() *db.Habit { return nil },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHabit := tt.setup()
			habit, err := GetBySlug(tt.slug, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, habit)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, habit)
				assert.Equal(t, habit.Slug, expectedHabit.Slug)
				assert.Equal(t, tt.user.ID, habit.UserId)
			}
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		newName string
		user    *db.User
		setup   func() *db.Habit
		wantErr bool
	}{
		{
			name:    "valid rename",
			slug:    "test-habit",
			newName: "New Name",
			user:    testUser,
			setup: func() *db.Habit {
				habit, _ := Create(createDTO{Name: "Test Habit"}, testUser)
				return habit
			},
			wantErr: false,
		},
		{
			name:    "non-existing habit",
			slug:    "non-existing",
			newName: "New Name",
			user:    testUser,
			setup:   func() *db.Habit { return nil },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := Rename(tt.slug, tt.user, tt.newName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				habit, _ := GetBySlug("new-name", tt.user)
				assert.NotNil(t, habit)
				assert.Equal(t, tt.newName, habit.Name)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		user    *db.User
		setup   func() *db.Habit
		wantErr bool
	}{
		{
			name: "existing habit",
			slug: "test-habit",
			user: testUser,
			setup: func() *db.Habit {
				habit, _ := Create(createDTO{Name: "Test Habit"}, testUser)
				return habit
			},
			wantErr: false,
		},
		{
			name:    "non-existing habit",
			slug:    "non-existing",
			user:    testUser,
			setup:   func() *db.Habit { return nil },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := Delete(tt.slug, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				habit, gserr := GetBySlug(tt.slug, tt.user)
				assert.NoError(t, gserr)
				assert.NotNil(t, habit)
			}
		})
	}
}

func cleanupTestData(t *testing.T) {
	db.Client.Exec("DELETE FROM habits")
}

func TestMain(m *testing.M) {
	utils.SetupTestEnv()
	m.Run()
	cleanupTestData(nil)
}
