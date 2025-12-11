package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mnuddindev/jutsu-api/internal/domain"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token has expired")
	ErrWrongTokenType  = errors.New("wrong token type")
	ErrTokenGeneration = errors.New("token generation failed")
)

type TokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
}

// Claims represents JWT token claims
type Claims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	AccessTTL    int
	RefreshTTL   int
}

// NewTokenService creates a new token service
func NewTokenService(accessSecret, refreshSecret, issuer string, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		issuer:        issuer,
	}
}

// GenerateTokenPair generates an access token and refresh token pair
func (ts *TokenService) GenerateTokenPair(user *domain.User) (*TokenPair, error) {
	accessToken, err := ts.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := ts.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessTTL:    int(ts.accessTTL.Seconds()),
		RefreshTTL:   int(ts.refreshTTL.Seconds()),
	}, nil
}

// GenerateAccessToken generates an access token
func (ts *TokenService) GenerateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   string(domain.AccessTokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ts.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    ts.issuer,
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.accessSecret)
}

// GenerateRefreshToken generates a refresh token
func (ts *TokenService) GenerateRefreshToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   string(domain.RefreshTokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ts.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    ts.issuer,
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.refreshSecret)
}

// ValidateAccessToken validates an access token and returns claims
func (ts *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return ts.validateToken(tokenString, ts.accessSecret, string(domain.AccessTokenType))
}

// ValidateRefreshToken validates a refresh token and returns claims
func (ts *TokenService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return ts.validateToken(tokenString, ts.refreshSecret, string(domain.RefreshTokenType))
}

// validateToken validates a token with given secret and expected type
func (ts *TokenService) validateToken(tokenString string, secret []byte, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != expectedType {
		return nil, ErrWrongTokenType
	}

	if claims.Issuer != ts.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// HashToken creates a SHA256 hash of the token (for storage)
func (ts *TokenService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ExtractTokenID extracts the JTI (token ID) from a token without validating
func (ts *TokenService) ExtractTokenID(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.ID, nil
}

// GetAccessTTL returns access token TTL
func (ts *TokenService) GetAccessTTL() time.Duration {
	return ts.accessTTL
}

// GetRefreshTTL returns refresh token TTL
func (ts *TokenService) GetRefreshTTL() time.Duration {
	return ts.refreshTTL
}
