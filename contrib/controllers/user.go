package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/contrib/models"
	"backend/contrib/repository"
	"backend/contrib/schema"
	"backend/contrib/slugify"
	"backend/core"
	"backend/validator"
)

// RetrieveUsersApi godoc
// @Summary Retrieves all users with pagination
// @Description Retrieves a list of users with pagination support
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search term"
// @Param is_active query bool false "User activity"
// @Success 200 {object} schema.UserPaginatedResponse
// @Failure 401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router /api/v1/user/ [get]
// @Security     OAuth2Password[]
func RetrieveUsersApi(c *fiber.Ctx, db *gorm.DB, commons core.CommonsModel) error {
	// Parse optional search
	search := c.Query("search", "")

	// Parse optional is_active
	var isActive *bool
	if v := c.Query("is_active"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isActive = &b
		}
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	userRepo := repository.NewUserRepository(db)
	users, err := userRepo.UserList(&repository.UserRepoFilter{
		Limit:          &commons.Limit,
		Offset:         &commons.Offset,
		Search:         &search,
		IsActive:       isActive, // added here
		IncludeDeleted: &userSession.User.IsSuperuser,
	})

	if err != nil {
		return fiber.ErrInternalServerError
	}
	userResponseList := schema.ToUserResponseList(users)

	// Return response
	return c.Status(200).JSON(schema.UserPaginatedResponse{
		Rows:  userResponseList,
		Page:  commons.Page,
		Limit: commons.Limit,
	})
}

// CountUsersApi godoc
// @Summary Count all users
// @Description Retrieves a count of users
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param is_active query bool false "User activity"
// @Success 200 {object} core.CountResponse "count: 12345"
// @Failure 401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router /api/v1/user/count/ [get]
// @Security     OAuth2Password[]
func CountUsersApi(c *fiber.Ctx, db *gorm.DB) error {
	// Parse optional search
	search := c.Query("search", "")

	// Parse optional is_active
	var isActive *bool
	if v := c.Query("is_active"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isActive = &b
		}
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	userRepo := repository.NewUserRepository(db)
	count, err := userRepo.UserCount(&repository.UserRepoFilter{
		Search:         &search,
		IsActive:       isActive, // added here
		IncludeDeleted: &userSession.User.IsSuperuser,
	})

	if err != nil {
		return fiber.ErrInternalServerError
	}

	// Return response
	return c.Status(200).JSON(core.CountResponse{
		Count: count,
	})
}

// CreateUserApi 	godoc
// @Summary 		Create and response new user
// @Description 	Creates new user
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param   		requestBody body schema.UserCreateInput true "User create Data"
// @Success 		201 {object} schema.UserMessageResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/user/create/ [post]
// @Security     	OAuth2Password[]
func CreateUserApi(c *fiber.Ctx, db *gorm.DB) error {
	// Parse and validate JSON body using ParseBody
	objIn, err := validator.ParseBody[schema.UserCreateInput](c)
	if err != nil {
		return err
	}
	if objIn.Role == models.OperatorRole && objIn.CarPark == nil {
		return &core.ValidationErrorResponse{
			Detail: "Validation error",
			Errors: map[string]string{
				"carPark": "Park number must for operator role",
			},
		}
	}

	userRepo := repository.NewUserRepository(db)
	includeSuperuser := true
	includeDeleted := true
	username := slugify.Make(objIn.Username)

	isExist, err := userRepo.UserExist(&repository.UserRepoFilter{
		Username:         &username,
		IncludeSuperuser: &includeSuperuser,
		IncludeDeleted:   &includeDeleted,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if isExist {
		return &core.ValidationErrorResponse{
			Detail: "Validation error",
			Errors: map[string]string{
				"username": "User with this username already exists",
			},
		}
	}
	user, err := userRepo.UserCreate(
		username,
		objIn.Password,
		objIn.FullName,
		objIn.Role,
		objIn.IsActive,
		objIn.CarPark,
		nil,
	)

	if err != nil {
		return fiber.ErrInternalServerError
	}
	userResponse := schema.ToUserResponse(*user)

	return c.Status(201).JSON(schema.UserMessageResponse{
		Data:    &userResponse,
		Message: "User created",
	})
}

// UserDetailApi 	godoc
// @Summary 		Retrieve single user
// @Description 	Retrieve single user detail
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID (UUID)"
// @Success 		200 {object} schema.UserResponse
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/user/{id}/detail/ [get]
// @Security     	OAuth2Password[]
func UserDetailApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	// Parse string ID to uuid.UUID
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.UserGetByID(userID, !userSession.User.IsSuperuser)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrNotFound
	}
	userResponse := schema.ToUserResponse(*user)
	return c.Status(200).JSON(userResponse)
}

// UpdateUserApi 	godoc
// @Summary 		Update a user
// @Description 	Updates user data by ID
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID (UUID)"
// @Param   		requestBody body schema.UserUpdateInput true "User update Data"
// @Success 		200 {object} schema.UserMessageResponse "Bad Request"
// @Failure 		400 {object} core.BadRequestResponse "detail: BadRequest - invalid request"
// @Failure 		401 {object} core.UnauthorizedResponse "detail: Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "detail: Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "detail: User not found"
// @Failure 		422 {object} core.ValidationErrorResponse "detail: Validation errors"
// @Failure 		500 {object} core.InternalServerErrorResponse "detail: Internal Server Error"
// @Router 			/api/v1/user/{id}/update/ [patch]
// @Security    	OAuth2Password[]
func UpdateUserApi(c *fiber.Ctx, db *gorm.DB) error {
	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}

	idStr := c.Params("id")
	// Parse string ID to uuid.UUID
	userId, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}
	if userId == userSession.UserID {
		return fiber.NewError(fiber.StatusBadRequest, "Can not update your self")
	}

	// Parse and validate JSON body using ParseBody
	objIn, err := validator.ParseBody[schema.UserUpdateInput](c)
	if err != nil {
		return err
	}

	if objIn.Role != nil && *objIn.Role == models.OperatorRole && objIn.CarPark == nil {
		return &core.ValidationErrorResponse{
			Detail: "Validation error",
			Errors: map[string]string{
				"carPark": "Park number must for operator role",
			},
		}
	}

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.UserGetByID(userId, !userSession.User.IsSuperuser)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrNotFound
	}
	includeSuperuser := true
	includeDeleted := true
	if objIn.Username != nil && *objIn.Username != "" && *objIn.Username != user.Username {
		username := slugify.Make(*objIn.Username)
		isExist, err := userRepo.UserExist(&repository.UserRepoFilter{
			Username:         &username,
			ExcludeID:        &userId,
			IncludeSuperuser: &includeSuperuser,
			IncludeDeleted:   &includeDeleted,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if isExist {
			return &core.ValidationErrorResponse{
				Detail: "Validation error",
				Errors: map[string]string{
					"username": "User with this username already exists",
				},
			}
		}
		user.Username = *objIn.Username
	}
	if objIn.Password != nil && *objIn.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*objIn.Password), bcrypt.DefaultCost)

		if err != nil {
			return fiber.ErrInternalServerError
		}
		user.Password = string(hashedPassword)
	}
	if objIn.FullName != nil && *objIn.FullName != "" {
		user.FullName = *objIn.FullName
	}
	if objIn.IsActive != nil {
		user.IsActive = *objIn.IsActive
	}
	if objIn.Role != nil {
		user.Role = *objIn.Role
	}
	if objIn.CarPark != nil {
		user.CarPark = objIn.CarPark
	}

	// Save update
	if err := userRepo.UserUpdate(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	userResponse := schema.ToUserResponse(*user)

	return c.Status(200).JSON(schema.UserMessageResponse{
		Data:    &userResponse,
		Message: "User updated",
	})
}

// UserDeleteApi 	godoc
// @Summary 		Delete a user
// @Description 	Deletes a user by ID
// @Tags 			Users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID (UUID)"
// @Success 		200 {object} schema.UserMessageResponse "User deleted successfully"
// @Failure 		400 {object} core.BadRequestResponse "Bad Request"
// @Failure 		401 {object} core.UnauthorizedResponse "Unauthorized - Invalid token"
// @Failure 		403 {object} core.PermissionDeniedResponse "Permission denied"
// @Failure 		404 {object} core.NotFoundResponse "User not found"
// @Failure 		500 {object} core.InternalServerErrorResponse "Internal Server Error"
// @Router 			/api/v1/user/{id}/delete/ [delete]
// @Security 		OAuth2Password[]
func UserDeleteApi(c *fiber.Ctx, db *gorm.DB) error {
	idStr := c.Params("id")
	// Parse string ID to uuid.UUID
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	us := c.Locals("userSession")
	if us == nil {
		return fiber.ErrUnauthorized
	}
	userSession, ok := us.(*models.UserSession)
	if !ok {
		return fiber.ErrInternalServerError
	}
	if userID == userSession.UserID {
		return fiber.NewError(fiber.StatusBadRequest, "Can not delete yourself")
	}

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.UserGetByID(userID, userSession.User.IsSuperuser)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrNotFound
	}
	deleteAt := time.Now()
	// Update username to prevent conflicts (add timestamp)
	user.Username = fmt.Sprintf("%s-deleted-%d", user.Username, deleteAt.Unix())
	// Mark as deleted
	user.DeletedAt = &gorm.DeletedAt{Time: deleteAt, Valid: true}
	user.IsActive = false

	if err := userRepo.UserUpdate(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(200).JSON(schema.UserMessageResponse{
		Data:    nil,
		Message: "User deleted",
	})
}
