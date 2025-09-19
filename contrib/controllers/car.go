package controllers

import (
	"backend/contrib/slugify"
	"backend/validator"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
)

// RetrieveCarsApi 	godoc
// @Summary 			Retrieves lists cars with optional search, type, and pagination
// @Description 		Retrieves a list of cars with pagination support
// @Tags 				Cars
// @Accept 				json
// @Produce 			json
// @Param 				page query int false "Page number"
// @Param 				limit query int false "Limit per page"
// @Param 				search query string false "Search term"
// @Success 			200 {object} schema.CarPaginatedResponse
// @Failure 			401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 			403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 			500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 				/api/v1/car/ [get]
// @Security     		OAuth2Password[]
func RetrieveCarsApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {
	search := c.Query("search", "")
	carRepo := repository.NewCarRepository(db)
	cars, err := carRepo.CarList(&repository.CarRepoFilter{
		Search: &search,
		Limit:  &commons.Limit,
		Offset: &commons.Offset,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	carResponseList := schema.ToCarResponseList(cars)
	return c.Status(200).JSON(schema.CarPaginatedResponse{
		Rows:  carResponseList,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// CountCarsApi 	godoc
// @Summary 		Count returns total cars count
// @Description 	Retrieves a count of cars
// @Tags 			Cars
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "Search term"
// @Success 		200 {object} core.CountResponse "count: 12345"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car/count/ [get]
// @Security     	OAuth2Password[]
func CountCarsApi(c *fiber.Ctx, db *gorm.DB) error {
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

// CreateCarApi 	godoc
// @Summary 		Creates a new car
// @Description 	Creates new car
// @Tags 			Cars
// @Accept 			json
// @Produce 		json
// @Param   		requestBody body schema.CarCreateInput true "Request body"
// @Success 		201 {object} schema.CarMessageResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car/create/ [post]
// @Security     	OAuth2Password[]
func CreateCarApi(c *fiber.Ctx, db *gorm.DB) error {
	objIn, err := validator.ParseBody[schema.CarCreateInput](c)
	if err != nil {
		return err
	}

	carRepo := repository.NewCarRepository(db)

	carNumber := slugify.Make(objIn.CarNumber)

	isExist, err := carRepo.CarExist(&repository.CarRepoFilter{
		CarNumber: &carNumber,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if isExist {
		return &core.ValidationErrorResponse{
			Detail: "Validation error",
			Errors: map[string]string{
				"carNumber": "Car with this car number already exists",
			},
		}
	}

	car, err := carRepo.CarCreate(objIn.CarNumber, objIn.IsStaff, objIn.OwnerName)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	carResponse := schema.ToCarResponse(*car)
	return c.Status(201).JSON(schema.CarMessageResponse{
		Data:    &carResponse,
		Message: "Car created",
	})
}

// CarDetailApi 	godoc
// @Summary 		Retrieve single car
// @Description 	Retrieve single car detail
// @Tags 			Cars
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Car ID (int)"
// @Success 		200 {object} schema.CarResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car/{id}/detail/ [get]
// @Security     	OAuth2Password[]
func CarDetailApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid car id")
	}

	carRepo := repository.NewCarRepository(db)
	car, err := carRepo.CarGetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if car == nil {
		return fiber.ErrNotFound
	}

	carResponse := schema.ToCarResponse(*car)
	return c.Status(200).JSON(carResponse)
}

// UpdateCarApi 	godoc
// @Summary 		Update a car
// @Description 	Updates an existing car
// @Tags 			Cars
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Car ID (int)"
// @Param   		requestBody body schema.CarUpdateInput true "Request data"
// @Success 		200 {object} schema.CarMessageResponse "Bad Request"
// @Failure 		400 {object} core.BadRequestResponse "detail: BadRequest - invalid request"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "detail: User not found"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/car/{id}/update/ [patch]
// @Security    	OAuth2Password[]
func UpdateCarApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid car id")
	}

	objIn, err := validator.ParseBody[schema.CarUpdateInput](c)
	if err != nil {
		return err
	}

	carRepo := repository.NewCarRepository(db)
	dbObj, err := carRepo.CarGetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if dbObj == nil {
		return fiber.ErrNotFound
	}

	if objIn.CarNumber != nil && *objIn.CarNumber != "" && *objIn.CarNumber != dbObj.CarNumber {
		carNumber := slugify.Make(*objIn.CarNumber)
		isExist, err := carRepo.CarExist(&repository.CarRepoFilter{
			CarNumber: &carNumber,
			ExcludeID: &dbObj.ID,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if isExist {
			return &core.ValidationErrorResponse{
				Detail: "Validation error",
				Errors: map[string]string{
					"carNumber": "Car with this car number already exists",
				},
			}
		}
		dbObj.CarNumber = *objIn.CarNumber
	}
	if objIn.OwnerName != nil {
		dbObj.OwnerName = objIn.OwnerName
	}

	if err := carRepo.CarUpdate(dbObj); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result := schema.ToCarResponse(*dbObj)
	return c.Status(200).JSON(schema.CarMessageResponse{
		Data:    &result,
		Message: "Car updated",
	})
}
