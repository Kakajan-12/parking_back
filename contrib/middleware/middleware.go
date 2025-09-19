package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/contrib/models"
	"backend/contrib/repository"
	"backend/core"
)

/*
Package auth provides authentication middleware and helpers for Fiber applications
using JWT tokens and GORM-based user sessions.

AuthToken middleware verifies JWT tokens, validates user sessions, and enforces
permission checks such as active status and admin privileges.

Usage:

	app := fiber.New()
	db := InitDB() // GORM DB instance

	authGroup := app.Group("/api/v1/auth")

	// Routes without authentication
	authGroup.Post("/register", func(c *fiber.Ctx) error { return RegisterApi(c, db) })
	authGroup.Post("/login", func(c *fiber.Ctx) error { return LoginApi(c, db) })

	// Protected routes
	authGroup.Get("/me", AuthToken(db, false, true), MeApi)                // active users only
	authGroup.Delete("/users/:id", AuthToken(db, true, false), DeleteUserApi) // admins (active enforced automatically)
*/

// AuthToken returns a Fiber middleware that authenticates requests using a JWT token.
//
// Parameters:
//   - db: *gorm.DB, database instance to fetch user sessions
//   - requireAdmin: bool, if true, only superusers can access the route. Automatically enforces requireActive.
//   - requireActive: bool, if true, only active users can access the route.
//
// Behavior:
//  1. Extracts JWT token from the Authorization header or cookie.
//  2. Decodes the token and extracts the "jti" field as session UUID.
//  3. Retrieves the user session from the database. Rejects if invalid or revoked.
//  4. Checks permissions:
//     - requireActive enforces user.IsActive == true
//     - requireAdmin enforces user.IsSuperuser == true
//  5. Stores the user session in c.Locals("userSession") for downstream handlers.
//
// Returns:
//
//	fiber.Handler, a middleware function that can be used in Fiber route chains.
func AuthToken(db *gorm.DB, requiredRoles []models.RoleType, requireActive bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenPtr := getTokenFromRequest(c)
		if tokenPtr == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "No token provided: 11101")
		}
		token := *tokenPtr

		claims, err := core.JWTDecode(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Parse UUID
		jti, err := ParseJwtUUID(claims, "jti")
		if err != nil {
			log.Debug("Failed to database: %v", err)
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: 11102")
		}
		userSessionRepo := repository.NewUserSessionRepository(db)
		userSession, err := userSessionRepo.UserSessionGetByID(*jti, true)

		if userSession == nil || err != nil {
			print(err, "ERR: 12020")
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: 11103")
		} else if userSession.RevokedAt != nil && userSession.RevokedAt.Before(time.Now()) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: 11104")
		}
		if userSession.User.DeletedAt != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: 11105")
		}

		if userSession.User.IsActive && userSession.User.IsSuperuser {
			c.Locals("userSession", userSession)
			return c.Next()
		}

		// Enforce active status if required
		if requireActive && !userSession.User.IsActive {
			return fiber.NewError(fiber.StatusUnauthorized, "User is inactive")
		}

		// Check role permissions
		if len(requiredRoles) > 0 {
			hasRole := false
			for _, role := range requiredRoles {
				if userSession.User.Role == role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
			}
		}

		// Store in locals for downstream handlers
		c.Locals("userSession", userSession)
		return c.Next()
	}
}

// Helper: extract token from Authorization header first, then cookie
func getTokenFromRequest(c *fiber.Ctx) *string {
	// Check Authorization header first
	authHeader := c.Get("Authorization")
	//println(authHeader, "authHeader")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return &token
	}

	// Fallback to cookie
	if token := c.Cookies("jwt"); token != "" {
		return &token
	}

	// Neither found
	return nil
}

func ParseJwtUUID(claims jwt.MapClaims, fieldName string) (*uuid.UUID, error) {
	// Extract user_id from claims
	fieldNameVal, ok := claims[fieldName]

	if !ok {
		return nil, errors.New("token is invalid: 11106")
	}

	// Convert to string
	var idStr string
	switch v := fieldNameVal.(type) {
	case string:
		idStr = v
	case float64:
		// JWT numeric values are float64
		idStr = fmt.Sprintf("%.0f", v)
	default:
		return nil, errors.New("token is invalid: 11107")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid UUID in token")
	}
	println(idStr, id.String(), "parsejwtuuid")
	return &id, nil
}
