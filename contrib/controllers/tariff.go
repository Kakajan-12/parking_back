package controllers

import (
	"backend/config"
	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
	"backend/validator"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RetrieveTariffsApi 	godoc
// @Summary 			Retrieves lists tariffs with optional search, type, and pagination
// @Description 		Retrieves a list of tariffs with pagination support
// @Tags 				Tariffs
// @Accept 				json
// @Produce 			json
// @Param 				page query int false "Page number"
// @Param 				limit query int false "Limit per page"
// @Success 			200 {object} schema.TariffPaginatedResponse
// @Failure 			401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 			403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 			500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 				/api/v1/tariff/ [get]
// @Security     		OAuth2Password[]
func RetrieveTariffsApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {

	objRepo := repository.NewTariffRepository(db)
	rows, err := objRepo.TariffList(&repository.TariffRepoFilter{
		Limit:  &commons.Limit,
		Offset: &commons.Offset,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	result := schema.ToTariffResponseList(rows)
	return c.Status(200).JSON(schema.TariffPaginatedResponse{
		Rows:  result,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// CountTariffsApi 	godoc
// @Summary 		Count returns total tariffs count
// @Description 	Retrieves a count of tariffs
// @Tags 			Tariffs
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} core.CountResponse "count: 12345"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/tariff/count/ [get]
// @Security     	OAuth2Password[]
func CountTariffsApi(c *fiber.Ctx, db *gorm.DB) error {
	objRepo := repository.NewTariffRepository(db)
	count, err := objRepo.TariffCount(nil)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(200).JSON(core.CountResponse{
		Count: count,
	})
}

// CreateTariffApi 	godoc
// @Summary 		Creates a new car
// @Description 	Creates new car
// @Tags 			Tariffs
// @Accept 			json
// @Produce 		json
// @Param   		requestBody body schema.TariffCreateInput true "Request body"
// @Success 		201 {object} schema.TariffMessageResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/tariff/create/ [post]
// @Security     	OAuth2Password[]
func CreateTariffApi(c *fiber.Ctx, db *gorm.DB) error {
	objIn, err := validator.ParseBody[schema.TariffCreateInput](c)
	if err != nil {
		return err
	}

	objRepo := repository.NewTariffRepository(db)
	duration := objIn.Duration

	isExist, err := objRepo.TariffExist(&repository.TariffRepoFilter{
		Duration: &duration,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if isExist {
		return &core.ValidationErrorResponse{
			Detail: "Validation error",
			Errors: map[string]string{
				"duration": "Tariff with this duration number already exists",
			},
		}
	}

	dbObj, err := objRepo.TariffCreate(objIn.Name, objIn.Duration, objIn.IsActive, objIn.PriceAmount, config.AppConfig.DefaultCurrency)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	result := schema.ToTariffResponse(*dbObj)
	return c.Status(201).JSON(schema.TariffMessageResponse{
		Data:    &result,
		Message: "Tariff created",
	})
}

// TariffDetailApi 	godoc
// @Summary 		Retrieve single tariff
// @Description 	Retrieve single tariff detail
// @Tags 			Tariffs
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Tariff ID (int)"
// @Success 		200 {object} schema.TariffResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/tariff/{id}/detail/ [get]
// @Security     	OAuth2Password[]
func TariffDetailApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid tariff id")
	}

	objRepo := repository.NewTariffRepository(db)
	dbObj, err := objRepo.TariffGetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if dbObj == nil {
		return fiber.ErrNotFound
	}

	result := schema.ToTariffResponse(*dbObj)
	return c.Status(200).JSON(result)
}

// UpdateTariffApi 	godoc
// @Summary 		Update a tariff
// @Description 	Updates an existing tariff
// @Tags 			Tariffs
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Tariff ID (int)"
// @Param   		requestBody body schema.TariffUpdateInput true "Request data"
// @Success 		200 {object} schema.TariffMessageResponse "Bad Request"
// @Failure 		400 {object} core.BadRequestResponse "detail: BadRequest - invalid request"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "detail: User not found"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/tariff/{id}/update/ [patch]
// @Security    	OAuth2Password[]
func UpdateTariffApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid car id")
	}

	objIn, err := validator.ParseBody[schema.TariffUpdateInput](c)
	if err != nil {
		return err
	}

	objRepo := repository.NewTariffRepository(db)
	dbObj, err := objRepo.TariffGetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if dbObj == nil {
		return fiber.ErrNotFound
	}

	if objIn.Duration != nil && *objIn.Duration != dbObj.Duration {
		isExist, err := objRepo.TariffExist(&repository.TariffRepoFilter{
			Duration:  objIn.Duration,
			ExcludeID: &dbObj.ID,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if isExist {
			return &core.ValidationErrorResponse{
				Detail: "Validation error",
				Errors: map[string]string{
					"duration": "Tariff with this duration already exists",
				},
			}
		}
		dbObj.Duration = *objIn.Duration
	}
	if objIn.Name != nil && *objIn.Name != "" {
		dbObj.Name = *objIn.Name
	}
	if objIn.IsActive != nil && *objIn.IsActive != dbObj.IsActive {
		dbObj.IsActive = *objIn.IsActive
	}
	if objIn.PriceAmount != nil && *objIn.PriceAmount != dbObj.PriceAmount {
		dbObj.PriceAmount = *objIn.PriceAmount
	}

	if err := objRepo.TariffUpdate(dbObj); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result := schema.ToTariffResponse(*dbObj)
	return c.Status(200).JSON(schema.TariffMessageResponse{
		Data:    &result,
		Message: "Tariff updated",
	})
}
