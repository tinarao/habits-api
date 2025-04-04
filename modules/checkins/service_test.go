package checkins

import (
	"hbapi/internal/db"
	"hbapi/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    *db.User
		habitId uint
		wantErr bool
	}{
		{
			name:    "valid creation",
			user:    &db.User{ID: 1},
			habitId: 1,
			wantErr: false,
		},
		{
			name:    "nil user",
			user:    nil,
			habitId: 1,
			wantErr: true,
		},
		{
			name:    "zero habit id",
			user:    &db.User{ID: 1},
			habitId: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkin, err := Create(tt.user, tt.habitId)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, checkin)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, checkin)
				assert.Equal(t, tt.user.ID, checkin.UserId)
				assert.Equal(t, tt.habitId, checkin.HabitId)
			}
		})
	}
}

func TestGetByHabit(t *testing.T) {
	tests := []struct {
		name    string
		habitId uint
		wantErr bool
		setup   func()
	}{
		{
			name:    "existing habit",
			habitId: 1,
			wantErr: false,
			setup: func() {
				// Create test checkins
				user := &db.User{ID: 1}
				Create(user, 1)
				Create(user, 1)
			},
		},
		{
			name:    "non-existing habit",
			habitId: 999,
			wantErr: false,
			setup:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			checkins, err := GetByHabit(tt.habitId)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, checkins)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, checkins)
				for _, checkin := range checkins {
					assert.Equal(t, tt.habitId, checkin.HabitId)
				}
			}
		})
	}
}

func TestGetLatest(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func()
	}{
		{
			name:    "get recent checkins",
			wantErr: false,
			setup: func() {
				user := &db.User{ID: 1}
				Create(user, 1)
				Create(user, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			checkins, err := GetLatest()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, checkins)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, checkins)

				// Verify all checkins are within last 48 hours
				twoDaysAgo := time.Now().UTC().Add(-48 * time.Hour)
				for _, checkin := range checkins {
					assert.True(t, checkin.CreatedAt.After(twoDaysAgo))
				}
			}
		})
	}
}

func cleanupTestData(t *testing.T) {
	db.Client.Exec("DELETE FROM check_ins")
}

func TestMain(m *testing.M) {
	utils.SetupTestEnv()
	m.Run()
	cleanupTestData(nil)
}
