package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Username      string     `json:"username" db:"username"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	Role          string     `json:"role" db:"role"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

// UserProfile represents extended user information
type UserProfile struct {
	UserID      uuid.UUID              `json:"user_id" db:"user_id"`
	AvatarURL   string                 `json:"avatar_url" db:"avatar_url"`
	Bio         string                 `json:"bio" db:"bio"`
	Preferences map[string]interface{} `json:"preferences" db:"preferences"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// UserRegisterRequest represents user registration payload
type UserRegisterRequest struct {
	Email            string `json:"email" validate:"required,email,max=255"`
	Username         string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password         string `json:"password" validate:"required,min=8,max=100"`
	Confirm_Password string `json:"confirm_password" validate:"required,min=8,max=100"`
}

// UserLoginRequest represents user login payload
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest represents user profile update payload
type UserUpdateRequest struct {
	Email     string  `json:"email" db:"email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

// ChangePasswordRequest represents password change payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         *User  `json:"user"`
}

// RefreshTokenRequest represents token refresh payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// PublicUser returns user data safe for public display (no sensitive info)
func (u *User) PublicUser() *User {
	return &User{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		EmailVerified: u.EmailVerified,
		Role:          u.Role,
		CreatedAt:     u.CreatedAt,
		LastLoginAt:   u.LastLoginAt,
	}
}

// IsValidRole checks if role is valid
func IsValidRole(role string) bool {
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
		"mod":   true,
	}
	return validRoles[role]
}
