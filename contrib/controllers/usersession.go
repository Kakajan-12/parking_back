package controllers

import (
	"backend/contrib/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
)

// RetrieveUserSessionsApi godoc
// @Summary 	Retrieves all users with pagination
// @Description Retrieves a list of users with pagination support
// @Tags 		Users
// @Accept 		json
// @Produce 	json
// @Param 		page query int false "Page number"
// @Param 		limit query int false "Limit per page"
// @Param 		user_id query string false "User ID (UUID)"
// @Success 	200 {object} schema.UserSessionExtendedPaginatedResponse
// @Failure 	401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 	403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 	500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 		/api/v1/user-session/ [get]
// @Security	OAuth2Password[]
func RetrieveUserSessionsApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {
	userSessionRepo := repository.NewUserSessionRepository(db)
	userIdr := c.Query("user_id", "")

	var userId *uuid.UUID
	if userIdr != "" {
		parsed, err := core.ParseUUID(userIdr)
		if err != nil {
			return fiber.ErrBadRequest
		} else {
			userId = &parsed
		}
	}

	userSessions, err := userSessionRepo.UserSessionList(&repository.UserSessionRepoFilter{
		Limit:  &commons.Limit,
		Offset: &commons.Offset,
		UserID: userId,
	}, true)

	if err != nil {
		return fiber.ErrInternalServerError
	}
	userSessionResponseList := schema.ToUserSessionExtendedResponseList(userSessions)
	// Return response
	return c.Status(200).JSON(schema.UserSessionExtendedPaginatedResponse{
		Rows:  userSessionResponseList,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// UserSessionDetailApi godoc
// @Summary Retrieve single user session
// @Description Retrieve single user session
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User session ID (UUID)"
// @Success 200 {object} schema.UserSessionExtendedResponse
// @Failure 401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router /api/v1/user-session/{id}/detail/ [get]
// @Security OAuth2Password[]
func UserSessionDetailApi(c *fiber.Ctx, db *gorm.DB) error {
	// Parse string ID to uuid.UUID
	objId, err := core.ParseUUID(c.Params("id", ""))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid path id")
	}

	repo := repository.NewUserSessionRepository(db)
	dbObj, err := repo.UserSessionGetByID(objId, true)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if dbObj == nil {
		return fiber.ErrNotFound
	}
	result := schema.ToUserSessionExtendedResponse(*dbObj)
	return c.Status(200).JSON(result)
}

// UserSessionRevokeApi godoc
// @Summary 		Retrieves all user sessions with pagination
// @Description 	Retrieves a list of users with pagination support
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User session ID (UUID)"
// @Success 		200 {object} core.MessagedResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/user-session/{id}/revoke/ [get]
// @Security 		OAuth2Password[]
func UserSessionRevokeApi(c *fiber.Ctx, db *gorm.DB) error {
	// Parse string ID to uuid.UUID
	objId, err := core.ParseUUID(c.Params("id", ""))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid path id")
	}
	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}
	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if objId == userSession.ID {
		return fiber.NewError(fiber.StatusBadRequest, "Can not revoke current session")
	}

	repo := repository.NewUserSessionRepository(db)
	count, err := repo.UserSessionRevoke(objId)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if count > 0 {
		return fiber.ErrNotFound
	}

	return c.Status(200).JSON(core.MessagedResponse{
		Message: "User session revoked",
	})
}
