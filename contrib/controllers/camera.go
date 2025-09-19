package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/contrib/models"
	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/core"
	"backend/validator"
)

// RetrieveCamerasApi 	godoc
// @Summary 			Retrieves lists cameras with optional search, type, and pagination
// @Description 		Retrieves a list of cameras with pagination support
// @Tags 				Cameras
// @Accept 				json
// @Produce 			json
// @Param 				page query int false "Page number"
// @Param 				limit query int false "Limit per page"
// @Param 				search query string false "Search term"
// @Param 				type query models.CameraType false "Camera type"
// @Success 			200 {object} schema.CameraPaginatedResponse
// @Failure 			401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 			403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 			500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 				/api/v1/camera/ [get]
// @Security     		OAuth2Password[]
func RetrieveCamerasApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {
	search := c.Query("search", "")
	cameraTypeStr := c.Query("type", "")

	var cameraType *models.CameraType
	if cameraTypeStr != "" {
		ct := models.CameraType(cameraTypeStr)
		cameraType = &ct
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	cameraRepo := repository.NewCameraRepository(db)
	cameras, err := cameraRepo.CameraList(&repository.CameraRepoFilter{
		Search:         &search,
		Type:           cameraType,
		Limit:          &commons.Limit,
		Offset:         &commons.Offset,
		IncludeDeleted: &userSession.User.IsSuperuser,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	cameraResponseList := schema.ToCameraResponseList(cameras)
	return c.Status(200).JSON(schema.CameraPaginatedResponse{
		Rows:  cameraResponseList,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// CountCamerasApi 	godoc
// @Summary 		Count returns total cameras count
// @Description 	Retrieves a count of cameras
// @Tags 			Cameras
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "Search term"
// @Param 			type query models.CameraType false "Camera type"
// @Success 		200 {object} core.CountResponse "count: 12345"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/camera/count/ [get]
// @Security     	OAuth2Password[]
func CountCamerasApi(c *fiber.Ctx, db *gorm.DB) error {
	search := c.Query("search", "")
	cameraTypeStr := c.Query("type", "")

	var cameraType *models.CameraType
	if cameraTypeStr != "" {
		ct := models.CameraType(cameraTypeStr)
		cameraType = &ct
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	cameraRepo := repository.NewCameraRepository(db)
	count, err := cameraRepo.CameraCount(&repository.CameraRepoFilter{
		Search:         &search,
		Type:           cameraType,
		IncludeDeleted: &userSession.User.IsSuperuser,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(200).JSON(core.CountResponse{
		Count: count,
	})
}

// CreateCameraApi 	godoc
// @Summary 		Creates a new camera
// @Description 	Creates new camera
// @Tags 			Cameras
// @Accept 			json
// @Produce 		json
// @Param   		requestBody body schema.CameraCreateInput true "Request body"
// @Success 		201 {object} schema.CameraMessageResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/camera/create/ [post]
// @Security     	OAuth2Password[]
func CreateCameraApi(c *fiber.Ctx, db *gorm.DB) error {
	objIn, err := validator.ParseBody[schema.CameraCreateInput](c)
	if err != nil {
		return err
	}

	cameraRepo := repository.NewCameraRepository(db)
	camera, err := cameraRepo.CameraCreate(objIn.Name, objIn.Type, objIn.ChannelName, objIn.ChannelToken)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	cameraResponse := schema.ToCameraResponse(*camera)
	return c.Status(201).JSON(schema.CameraMessageResponse{
		Data:    &cameraResponse,
		Message: "Camera created",
	})
}

// CameraDetailApi 	godoc
// @Summary 		Retrieve single camera
// @Description 	Retrieve single camera detail
// @Tags 			Cameras
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Camera ID (int)"
// @Success 		200 {object} schema.CameraResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/camera/{id}/detail/ [get]
// @Security     	OAuth2Password[]
func CameraDetailApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid camera id")
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	cameraRepo := repository.NewCameraRepository(db)
	camera, err := cameraRepo.CameraGetByID(id, !userSession.User.IsSuperuser)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if camera == nil {
		return fiber.ErrNotFound
	}

	cameraResponse := schema.ToCameraResponse(*camera)
	return c.Status(200).JSON(cameraResponse)
}

// UpdateCameraApi 	godoc
// @Summary 		Update a camera
// @Description 	Updates an existing camera
// @Tags 			Cameras
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Camera ID (int)"
// @Param   		requestBody body schema.CameraUpdateInput true "Request data"
// @Success 		200 {object} schema.CameraMessageResponse "Bad Request"
// @Failure 		400 {object} core.BadRequestResponse "detail: BadRequest - invalid request"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "detail: User not found"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/camera/{id}/update/ [patch]
// @Security    	OAuth2Password[]
func UpdateCameraApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid camera id")
	}

	objIn, err := validator.ParseBody[schema.CameraUpdateInput](c)
	if err != nil {
		return err
	}

	cameraRepo := repository.NewCameraRepository(db)
	camera, err := cameraRepo.CameraGetByID(id, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if camera == nil {
		return fiber.ErrNotFound
	}

	if objIn.Name != nil {
		camera.Name = *objIn.Name
	}
	if objIn.Type != nil {
		camera.Type = *objIn.Type
	}

	if err := cameraRepo.CameraUpdate(camera); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cameraResponse := schema.ToCameraResponse(*camera)
	return c.Status(200).JSON(schema.CameraMessageResponse{
		Data:    &cameraResponse,
		Message: "Camera updated",
	})
}

// CameraDeleteApi 	godoc
// @Summary 		Deletes a camera
// @Description 	Deletes a camera by ID
// @Tags 			Cameras
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "Camera ID (int)"
// @Success 		200 {object} schema.UserMessageResponse "User deleted successfully"
// @Failure 		400 {object} core.BadRequestResponse "Bad Request"
// @Failure 		401 {object} core.UnauthorizedResponse "Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "User not found"
// @Failure 		500 {object} core.InternalServerErrorResponse "Internal Server Error"
// @Router 			/api/v1/camera/{id}/delete/ [delete]
// @Security 		OAuth2Password[]
func CameraDeleteApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid camera id")
	}

	cameraRepo := repository.NewCameraRepository(db)
	camera, err := cameraRepo.CameraGetByID(id, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if camera == nil {
		return fiber.ErrNotFound
	}

	deleteAt := time.Now()
	camera.Name = fmt.Sprintf("%s-deleted-%d", camera.Name, deleteAt.Unix())
	camera.DeletedAt = &gorm.DeletedAt{Time: deleteAt, Valid: true}

	if err := cameraRepo.CameraUpdate(camera); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(200).JSON(schema.CameraMessageResponse{
		Data:    nil,
		Message: "Camera deleted",
	})
}
