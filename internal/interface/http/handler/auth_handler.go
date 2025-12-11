package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/mnuddindev/jutsu-api/internal/domain"
	"github.com/mnuddindev/jutsu-api/internal/interface/validation"
	"github.com/mnuddindev/jutsu-api/internal/repo"
	"github.com/mnuddindev/jutsu-api/internal/service"
)

// AuthHandler handles authentication HTTP endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// getDeviceInfo extracts device information from request
func (h *AuthHandler) getDeviceInfo(c *fiber.Ctx) string {
	return c.Get("User-Agent", "Unknown")
}

// getIPAddress extracts IP address from request
func (h *AuthHandler) getIPAddress(c *fiber.Ctx) string {
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		if idx := strings.Index(ip, ","); idx > 0 {
			return strings.TrimSpace(ip[:idx])
		}
		return ip
	}
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return c.IP()
}

// getUserAgent extracts the User-Agent string from the request,
// checking multiple common proxy/CDN headers and normalizing the value.
func (h *AuthHandler) getUserAgent(c *fiber.Ctx) string {
	headers := []string{
		"User-Agent",
		"X-User-Agent",
		"X-Device-User-Agent",
		"X-Original-User-Agent",
		"X-OperaMini-Phone-UA",
		"X-UCBrowser-UA",
		"X-WAP-Profile",
	}

	for _, header := range headers {
		if ua := strings.TrimSpace(c.Get(header)); ua != "" {
			if idx := strings.Index(ua, ","); idx > 0 {
				return strings.TrimSpace(ua[:idx])
			}
			return ua
		}
	}

	return "unknown"
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateRequest(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	user, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, repo.ErrEmailExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		if errors.Is(err, repo.ErrUsernameExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": err.Error(),
			"error":  "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user.PublicUser(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateRequest(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	deviceInfo := h.getDeviceInfo(c)
	ipAddress := h.getIPAddress(c)
	userAgent := h.getUserAgent(c)

	authResp, err := h.authService.Login(c.Context(), &req, deviceInfo, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		if errors.Is(err, service.ErrUserNotActive) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User account is not active",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  err.Error(),
			"message": "Login failed",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResp.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/api/auth",
		MaxAge:   7 * 24 * 60 * 60,
	})

	return c.JSON(fiber.Map{
		"access_token": authResp.AccessToken,
		"token_type":   authResp.TokenType,
		"expires_in":   authResp.ExpiresIn,
		"user":         authResp.User,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		var req domain.RefreshTokenRequest
		if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Refresh token not found",
			})
		}
	}

	deviceInfo := h.getDeviceInfo(c)
	ipAddress := h.getIPAddress(c)
	userAgent := h.getUserAgent(c)

	authResp, err := h.authService.RefreshToken(c.Context(), refreshToken, deviceInfo, ipAddress, userAgent)
	if err != nil {
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HTTPOnly: true,
			Secure:   true,
			MaxAge:   -1,
		})

		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired refresh token",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResp.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/api/auth",
		MaxAge:   7 * 24 * 60 * 60,
	})

	return c.JSON(fiber.Map{
		"access_token": authResp.AccessToken,
		"token_type":   authResp.TokenType,
		"expires_in":   authResp.ExpiresIn,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var accessToken string
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	}

	refreshToken := c.Cookies("refresh_token")

	if err := h.authService.Logout(c.Context(), accessToken, refreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Logout failed",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// LogoutAll handles logout from all devices
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authHeader := c.Get("Authorization")
	var accessToken string
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	}

	if err := h.authService.LogoutAll(c.Context(), userID, accessToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout from all devices",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{
		"message": "Logged out from all devices successfully"})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	user, err := h.authService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user profile",
		})
	}

	sessions, _ := h.authService.GetActiveSessions(c.Context(), userID)

	return c.JSON(fiber.Map{
		"user":            user,
		"active_sessions": sessions,
	})
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var req domain.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateRequest(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	user, err := h.authService.UpdateUser(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, repo.ErrUsernameExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var req domain.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateRequest(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	err := h.authService.ChangePassword(c.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Current password is incorrect",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change password",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{
		"message": "Password changed successfully. Please login again.",
	})
}
