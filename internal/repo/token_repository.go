package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"go.uber.org/zap"
)

// TokenRepository handles refresh token storage in Redis and PostgreSQL
type TokenRepository struct {
	db     *pgxpool.Pool
	cache  *cache.Manager
	logger *zap.Logger
}

// RefreshTokenData represents refresh token metadata stored in Redis
type RefreshTokenData struct {
	UserID     uuid.UUID `json:"user_id"`
	TokenHash  string    `json:"token_hash"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *pgxpool.Pool, cacheManager *cache.Manager, logger *zap.Logger) *TokenRepository {
	return &TokenRepository{
		db:     db,
		cache:  cacheManager,
		logger: logger,
	}
}

// StoreRefreshToken stores a refresh token in both Redis (primary) and PostgreSQL (backup)
func (r *TokenRepository) StoreRefreshToken(
	ctx context.Context,
	tokenHash string,
	userID uuid.UUID,
	expiresAt time.Time,
	deviceInfo, ipAddress, userAgent string,
) error {
	tokenData := RefreshTokenData{
		UserID:     userID,
		TokenHash:  tokenHash,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
	}

	// Store in Redis using cache manager
	if err := r.cache.Set(ctx, cache.CategoryRefreshToken, tokenHash, tokenData); err != nil {
		r.logger.Error("failed to store token in redis",
			zap.String("token_hash", tokenHash),
			zap.Error(err),
		)
	}

	// Add to user's active sessions set
	userSessionsKey := userID.String()
	if err := r.cache.HSet(ctx, cache.CategoryUserSession, userSessionsKey, tokenHash, tokenHash); err != nil {
		r.logger.Warn("failed to add to user sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	// Handle IP properly for PostgreSQL
	var ip sql.NullString
	if ipAddress != "" {
		ip = sql.NullString{String: ipAddress, Valid: true}
	} else {
		ip = sql.NullString{Valid: false} // will store NULL in DB
	}
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (token_hash) DO UPDATE 
		SET expires_at = EXCLUDED.expires_at
	`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, deviceInfo, ip, expiresAt)
	if err != nil {
		r.logger.Error("failed to store token in PostgreSQL",
			zap.String("token_hash", tokenHash),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// GetRefreshToken retrieves refresh token data
func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenData, error) {
	// Try Redis first (faster)
	var tokenData RefreshTokenData
	data, err := r.cache.Get(ctx, cache.CategoryRefreshToken, tokenHash)

	if err == nil {
		if err := json.Unmarshal(data, &tokenData); err != nil {
			r.logger.Error("failed to unmarshal token data",
				zap.String("token_hash", tokenHash),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
		}
		return &tokenData, nil
	}

	// If not in Redis, try PostgreSQL backup
	r.logger.Debug("token not in cache, checking database",
		zap.String("token_hash", tokenHash),
	)
	return r.getRefreshTokenFromDB(ctx, tokenHash)
}

// getRefreshTokenFromDB retrieves token from PostgreSQL
func (r *TokenRepository) getRefreshTokenFromDB(ctx context.Context, tokenHash string) (*RefreshTokenData, error) {
	query := `
		SELECT user_id, token_hash, device_info, ip_address, created_at, expires_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > NOW() AND revoked_at IS NULL
	`

	var tokenData RefreshTokenData
	var deviceInfo, ipAddress *string

	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&tokenData.UserID,
		&tokenData.TokenHash,
		&deviceInfo,
		&ipAddress,
		&tokenData.CreatedAt,
		&tokenData.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("token not found or expired: %w", err)
	}

	if deviceInfo != nil {
		tokenData.DeviceInfo = *deviceInfo
	}
	if ipAddress != nil {
		tokenData.IPAddress = *ipAddress
	}

	return &tokenData, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	// Remove from Redis using cache manager
	if err := r.cache.Delete(ctx, cache.CategoryRefreshToken, tokenHash); err != nil {
		r.logger.Warn("failed to delete token from redis",
			zap.String("token_hash", tokenHash),
			zap.Error(err),
		)
	}

	// Mark as revoked in PostgreSQL
	query := `
		UPDATE refresh_tokens
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE token_hash = $1
	`
	_, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		r.logger.Error("failed to revoke token in database",
			zap.String("token_hash", tokenHash),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke token in database: %w", err)
	}

	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	// Delete user sessions using cache manager
	userSessionsKey := userID.String()
	if err := r.cache.Delete(ctx, cache.CategoryUserSession, userSessionsKey); err != nil {
		r.logger.Warn("failed to delete user sessions from cache",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	// Revoke all tokens in PostgreSQL
	query := `
		UPDATE refresh_tokens
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.Error("failed to revoke user tokens in database",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke user tokens in database: %w", err)
	}

	return nil
}

// BlacklistAccessToken blacklists an access token (for logout)
func (r *TokenRepository) BlacklistAccessToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	// Use SetString with custom TTL for blacklisted tokens
	if err := r.cache.SetString(ctx, cache.CategoryJWTBlacklist, tokenID, "revoked"); err != nil {
		r.logger.Error("failed to blacklist access token",
			zap.String("token_id", tokenID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to blacklist access token: %w", err)
	}

	// Set custom expiration matching token TTL
	if err := r.cache.Expire(ctx, cache.CategoryJWTBlacklist, tokenID, ttl); err != nil {
		r.logger.Warn("failed to set expiration on blacklisted token",
			zap.String("token_id", tokenID),
			zap.Error(err),
		)
	}

	return nil
}

// IsAccessTokenBlacklisted checks if an access token is blacklisted
func (r *TokenRepository) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	exists, err := r.cache.Exists(ctx, cache.CategoryJWTBlacklist, tokenID)
	if err != nil {
		r.logger.Error("failed to check token blacklist",
			zap.String("token_id", tokenID),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return exists, nil
}

// GetUserActiveSessions returns count of active sessions for a user
func (r *TokenRepository) GetUserActiveSessions(ctx context.Context, userID uuid.UUID) (int, error) {
	// Query PostgreSQL for active sessions count
	query := `
		SELECT COUNT(*) 
		FROM refresh_tokens 
		WHERE user_id = $1 
		AND expires_at > NOW() 
		AND revoked_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		r.logger.Error("failed to get active sessions count",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to get active sessions: %w", err)
	}

	return count, nil
}

// CleanupExpiredTokens removes expired tokens from PostgreSQL (run as cron job)
func (r *TokenRepository) CleanupExpiredTokens(ctx context.Context) (int, error) {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW() OR revoked_at IS NOT NULL
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return int(result.RowsAffected()), nil
}
