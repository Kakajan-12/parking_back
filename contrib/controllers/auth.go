package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/config"
	"backend/contrib/models"
	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
	"backend/validator"
)

// LoginApi godoc
// @Summary      User login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        requestBody body schema.LoginInput true "User Login Data"
// @Success      200 {object} schema.LoginResponse "Login successful"
// @Failure      400 {object} core.BadRequestResponse "Invalid request body"
// @Failure      401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure      422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure      500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router       /api/v1/auth/login/ [post]
func LoginApi(c *fiber.Ctx, db *gorm.DB) error {
	loginInput, err := validator.ParseBody[schema.LoginInput](c)

	if err != nil {
		return err
	}

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.UserAuthenticate(loginInput.Username, loginInput.Password)

	if err != nil || user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	if !user.IsActive {
		return fiber.NewError(fiber.StatusUnauthorized, "Account is not active")
	}

	userSessionRepo := repository.NewUserSessionRepository(db)

	iat := time.Now()
	exp := iat.Add(config.AppConfig.JwtExpiration)

	ips := c.IPs()
	for i, ip := range ips {
		if strings.HasPrefix(ip, "::ffff:") {
			ips[i] = strings.TrimPrefix(ip, "::ffff:")
		}
	}
	ipStr := strings.Join(ips, "|")
	fmt.Println("ips slice:", ips) // prints the actual slice
	fmt.Println("ipStr:", ipStr)   // prints the joined string

	userAgent := c.Get("User-Agent")
	userSession, err := userSessionRepo.UserSessionCreate(
		user.ID,
		exp,
		&ipStr,
		&userAgent,
	)

	if err != nil {
		return fiber.ErrInternalServerError
	}
	_, err = userSessionRepo.UserSessionRevokeAllByUserID(user.ID, &userSession.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	payload := map[string]interface{}{
		"user_id": user.ID,
		"role":    "admin",
		"jti":     userSession.ID,
	}
	claims := core.JWTPayload(payload, &iat, &config.AppConfig.JwtExpiration)
	token, err := core.JWTEncode(claims)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(schema.LoginResponse{
		AccessToken: token,
		Role:        user.Role,
	})
}

// LogoutApi godoc
// @Summary Logout User
// @Description  Ends the session of a logged-in user by deleting the JWT token cookie.
// @Tags         Auth
// @Produce      json
// @Success      200 {object} core.MessagedResponse "message: Logout successful"
// @Failure      401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure      500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router       /api/v1/auth/logout/ [get]
// @Security     OAuth2Password[]
func LogoutApi(c *fiber.Ctx, db *gorm.DB) error {
	// Retrieve userSession from locals
	sess := c.Locals("userSession")
	userSession, ok := sess.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}
	now := time.Now()
	userSession.RevokedAt = &now

	userSessionRepo := repository.NewUserSessionRepository(db)
	err := userSessionRepo.UserSessionUpdate(userSession)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// MeApi @Summary Get current user information
// @Description  Retrieves the current user's username, role, and user ID from the JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} schema.UserSessionExtendedResponse  "detail: Returns user session information"
// @Failure      401 {object} core.ErrorResponse "detail: Unauthorized - Invalid token"
// @Failure      500 {object} core.ErrorResponse "detail: Internal Server Error"
// @Router       /api/v1/auth/me/ [get]
// @Security     OAuth2Password[]
func MeApi(c *fiber.Ctx) error {
	// Retrieve userSession from locals
	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}
	// Type assert to *UserSession
	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	// Convert to DTO
	response := schema.ToUserSessionExtendedResponse(*userSession)
	return c.JSON(response)
}
