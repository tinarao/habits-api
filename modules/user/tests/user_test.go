package tests

import (
	"fmt"
	"hbapi/internal/db"
	"hbapi/modules/auth"
	"hbapi/modules/user"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
)

func setupTestEnv() error {
	envFilepath := "path-to-.env"
	fmt.Println(envFilepath)

	if err := godotenv.Load(envFilepath); err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}

	db.Init()
	auth.Init()

	return nil
}

func TestCompleteAuth(t *testing.T) {
	tests := []struct {
		name     string
		gothUser *goth.User
		wantErr  bool
	}{
		{
			name: "new user auth",
			gothUser: &goth.User{
				Email:        "new@example.com",
				Name:         "New User",
				Provider:     "google",
				AvatarURL:    "https://example.com/avatar.jpg",
				RefreshToken: "refresh-token",
				ExpiresAt:    time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
	}

	if err := setupTestEnv(); err != nil {
		t.Fatalf("failed to setup test environment: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.CompleteAuth(tt.gothUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompleteAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.gothUser.Email, got.Email)
			}
		})
	}
}

func TestFindByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantExist bool
	}{
		{
			name:      "non-existing email",
			email:     "nonexistent@example.com",
			wantExist: false,
		},
		{
			name:      "empty email",
			email:     "",
			wantExist: false,
		},
	}

	if err := setupTestEnv(); err != nil {
		t.Fatalf("failed to setup test environment: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, exists := user.FindByEmail(tt.email)
			assert.Equal(t, tt.wantExist, exists)
			if tt.wantExist {
				assert.NotNil(t, got)
				assert.Equal(t, tt.email, got.Email)
			}
		})
	}
}

func TestPersistAccount(t *testing.T) {
	name := "Test User"
	avatarURL := "https://example.com/avatar.jpg"
	refreshToken := "refresh-token"

	tests := []struct {
		name     string
		gothUser *goth.User
		wantErr  bool
	}{
		{
			name: "successful persistence",
			gothUser: &goth.User{
				Email:        "test@example.com",
				Name:         name,
				Provider:     "google",
				AvatarURL:    avatarURL,
				RefreshToken: refreshToken,
				ExpiresAt:    time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "missing email",
			gothUser: &goth.User{
				Name:     "Test User",
				Provider: "google",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.PersistAccount(tt.gothUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersistAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.gothUser.Email, got.Email)
				assert.Equal(t, tt.gothUser.Name, *got.Name)
				assert.Equal(t, tt.gothUser.Provider, got.Provider)
				assert.Equal(t, tt.gothUser.AvatarURL, *got.ImageUrl)
				assert.Equal(t, tt.gothUser.RefreshToken, *got.RefreshToken)
				assert.Equal(t, tt.gothUser.ExpiresAt, got.RefreshTokenExpires)
				assert.Equal(t, db.FreePlan, got.Plan)
			}
		})
	}
}
