package controllers

import (
	"backend/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
)

// RetrieveCarSessionsApi 	godoc
// @Tags 				Cars
// @Accept 				json
// @Produce 			json
// @Param 				page query int false "Page number"
// @Param 				limit query int false "Limit per page"
// @Success 			200 {object} schema.CarSessionPaginatedResponse
// @Failure 			401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 			403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 			500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 				/api/v1/car-session/ [get]
// @Security     		OAuth2Password[]
func RetrieveCarSessionsApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {
	search := c.Query("search", "")
	objRepo := repository.NewCarSessionRepository(db)
	rows, err := objRepo.CarSessionList(&repository.CarSessionRepoFilter{
		Search: &search,
		Limit:  &commons.Limit,
		Offset: &commons.Offset,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	result := schema.ToCarSessionResponseList(rows)
	return c.Status(200).JSON(schema.CarSessionPaginatedResponse{
		Rows:  result,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// CountCarSessionsApi 	godoc
// @Tags 			Cars
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} core.CountResponse "count: 12345"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car-session/count/ [get]
// @Security     	OAuth2Password[]
func CountCarSessionsApi(c *fiber.Ctx, db *gorm.DB) error {
	search := c.Query("search", "")
	carRepo := repository.NewCarRepository(db)
	count, err := carRepo.CarCount(&repository.CarRepoFilter{
		Search: &search,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(200).JSON(core.CountResponse{
		Count: count,
	})
}

// CarSessionEventEntryApi 	godoc
// @Tags 			Cars
// @Param   		requestBody body map[string]interface{} true "Request body"
// @Success 		200 {object} schema.CarSessionMessageResponse "count: 12345"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car-session/event/entry/ [post]
func CarSessionEventEntryApi(c *fiber.Ctx, db *gorm.DB) error {
	objIn, err := validator.ParseBody[map[string]interface{}](c)
	if err != nil {
		return err
	}
	println(objIn)
	return c.Status(201).JSON(schema.CarSessionMessageResponse{
		Message: "Event entry created",
		Data:    nil,
	})
}
