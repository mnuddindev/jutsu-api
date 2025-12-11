package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenType represents the type of token
type TokenType string

const (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
)

// RefreshToken represents a stored refresh token
type RefreshToken struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash  string     `json:"-" db:"token_hash"` // Never expose
	DeviceInfo string     `json:"device_info,omitempty" db:"device_info"`
	IPAddress  string     `json:"ip_address,omitempty" db:"ip_address"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// IsValid checks if the refresh token is still valid
func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// TokenClaims represents the claims stored in JWT
type TokenClaims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	JTI    string    `json:"jti"` // JWT ID
}

// AccessTokenResponse represents access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// BlacklistedToken represents a blacklisted access token
type BlacklistedToken struct {
	TokenID   string    `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// SessionInfo represents active session information
type SessionInfo struct {
	TokenHash  string    `json:"token_hash"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsCurrent  bool      `json:"is_current"`
}
