package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mnuddindev/jutsu-api/internal/domain"
	"github.com/mnuddindev/jutsu-api/internal/infrastructure/logger"
	"github.com/mnuddindev/jutsu-api/internal/repo"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrEmailNotVerified   = errors.New("email not verified")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  *repo.UserRepository
	tokenRepo *repo.TokenRepository
	tokenSvc  *TokenService
	passSvc   *PasswordService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repo.UserRepository,
	tokenRepo *repo.TokenRepository,
	tokenSvc *TokenService,
	passSvc *PasswordService,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		tokenSvc:  tokenSvc,
		passSvc:   passSvc,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *domain.UserRegisterRequest) (*domain.User, error) {
	// Check if email exists
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, repo.ErrEmailExists
	}

	// Check if username exists
	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if usernameExists {
		return nil, repo.ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := s.passSvc.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:         req.Email,
		Username:      req.Username,
		PasswordHash:  string(hashedPassword),
		Role:          "user", // Default role
		EmailVerified: false,
		IsActive:      true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Send verification email
	// s.sendVerificationEmail(ctx, user)

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(
	ctx context.Context,
	req *domain.UserLoginRequest,
	deviceInfo, ipAddress, userAgent string,
) (*domain.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := s.passSvc.ComparePass(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	activeSessions, err := s.tokenRepo.GetUserActiveSessions(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not verify login status")
	}

	if activeSessions >= 1 {
		return nil, fiber.NewError(fiber.StatusForbidden, "Maximum login sessions reached")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Generate token pair
	tokenPair, err := s.tokenSvc.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshTokenHash := s.tokenSvc.HashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(s.tokenSvc.GetRefreshTTL())

	if err := s.tokenRepo.StoreRefreshToken(
		ctx,
		refreshTokenHash,
		user.ID,
		expiresAt,
		deviceInfo,
		ipAddress,
		userAgent,
	); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		logger.Warn("Failed to update last login",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
	}

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.AccessTTL,
		User:         user.PublicUser(),
	}, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
	deviceInfo, ipAddress, userAgent string,
) (*domain.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.tokenSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if token exists in storage
	tokenHash := s.tokenSvc.HashToken(refreshToken)
	tokenData, err := s.tokenRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired: %w", err)
	}

	// Verify user ID matches
	if tokenData.UserID != claims.UserID {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Revoke old refresh token (token rotation)
	if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		logger.Warn("Failed to revoke old refresh token",
			zap.String("user_id", claims.UserID.String()),
			zap.Error(err),
		)
	}

	// Generate new token pair
	tokenPair, err := s.tokenSvc.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store new refresh token
	newRefreshTokenHash := s.tokenSvc.HashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(s.tokenSvc.GetRefreshTTL())

	if err := s.tokenRepo.StoreRefreshToken(
		ctx,
		newRefreshTokenHash,
		user.ID,
		expiresAt,
		deviceInfo,
		ipAddress,
		userAgent,
	); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.AccessTTL,
		User:         user.PublicUser(),
	}, nil
}

// Logout revokes user's current session tokens
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// Blacklist access token
	if accessToken != "" {
		tokenID, err := s.tokenSvc.ExtractTokenID(accessToken)
		if err == nil {
			ttl := s.tokenSvc.GetAccessTTL()
			if err := s.tokenRepo.BlacklistAccessToken(ctx, tokenID, ttl); err != nil {
				logger.Warn("Failed to blacklist access token",
					zap.String("token_id", tokenID),
					zap.Error(err),
				)
			}
		}
	}

	// Revoke refresh token
	if refreshToken != "" {
		tokenHash := s.tokenSvc.HashToken(refreshToken)
		if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			logger.Warn("Failed to revoke refresh token", zap.Error(err))
		}
	}

	return nil
}

// LogoutAll revokes all user sessions
func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID, currentAccessToken string) error {
	// Blacklist current access token
	if currentAccessToken != "" {
		tokenID, err := s.tokenSvc.ExtractTokenID(currentAccessToken)
		if err == nil {
			ttl := s.tokenSvc.GetAccessTTL()
			if err := s.tokenRepo.BlacklistAccessToken(ctx, tokenID, ttl); err != nil {
				logger.Warn("Failed to blacklist access token",
					zap.String("token_id", tokenID),
					zap.Error(err),
				)
			}
		}
	}

	// Revoke all refresh tokens
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}

	return nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	currentPassword, newPassword string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := s.passSvc.ComparePass(currentPassword, user.PasswordHash); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := s.passSvc.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all tokens for security
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		logger.Warn("Failed to revoke all tokens after password change",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	return nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.PublicUser(), nil
}

// UpdateUser updates user profile
func (s *AuthService) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UserUpdateRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Username != nil {
		// Check if username is taken by another user
		if *req.Username != user.Username {
			exists, err := s.userRepo.UsernameExists(ctx, *req.Username)
			if err != nil {
				return nil, fmt.Errorf("failed to check username: %w", err)
			}
			if exists {
				return nil, repo.ErrUsernameExists
			}
			user.Username = *req.Username
		}
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user.PublicUser(), nil
}

// VerifyEmail marks user email as verified
func (s *AuthService) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.EmailVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// GetActiveSessions returns count of active sessions
func (s *AuthService) GetActiveSessions(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.tokenRepo.GetUserActiveSessions(ctx, userID)
}
