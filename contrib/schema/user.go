package schema

import (
	"time"

	"backend/contrib/models"
	"backend/core"
)

var roleLabels = map[models.RoleType]string{
	models.AdminRole:      "Administrator",
	models.OperatorRole:   "Operator",
	models.AccountantRole: "Accountant",
}

func RoleToChoice(role models.RoleType) core.ChoiceBase {
	return core.ChoiceBase{
		Value: string(role),
		Label: roleLabels[role],
	}
}

type UserUpdateInput struct {
	Username *string             `json:"username,omitempty" validate:"omitempty,max=150"`
	Password *string             `json:"password,omitempty" validate:"omitempty,min=8,max=128" example:"secret-word"`
	FullName *string             `json:"fullName,omitempty" validate:"omitempty,max=255"`
	IsActive *bool               `json:"isActive,omitempty" validate:"omitempty"`
	Role     *models.RoleType    `json:"role,omitempty"     validate:"omitempty,roletype" example:"operator"`
	CarPark  *models.CarParkType `json:"carPark,omitempty"  validate:"omitempty,carparktype"` // optional
} // @name UserUpdateInput

type UserCreateInput struct {
	Username string              `json:"username" validate:"required,max=150"`
	Password string              `json:"password" validate:"required,min=8,max=128" example:"secret-word"`
	FullName string              `json:"fullName" validate:"required,max=255"`
	IsActive bool                `json:"isActive"`
	Role     models.RoleType     `json:"role" validate:"required,roletype" example:"operator"`
	CarPark  *models.CarParkType `json:"carPark" validate:"omitempty,carparktype"` // optional
} // @name UserCreateInput

type UserResponse struct {
	ID        string           `json:"id" validate:"required"`
	Username  string           `json:"username" validate:"required"`
	FullName  string           `json:"fullName"` // required in TS, can be empty
	IsActive  bool             `json:"isActive" validate:"required"`
	Role      core.ChoiceBase  `json:"role" validate:"required"`
	CarPark   *core.ChoiceBase `json:"carPark"` // optional
	CreatedAt time.Time        `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time       `json:"updatedAt"` // optional, can be nil
} // @name UserResponse

func ToUserResponse(value models.User) UserResponse {
	var carPark *core.ChoiceBase
	if value.CarPark != nil {
		c := CarParkToChoice(string(*value.CarPark))
		carPark = &c
	}

	return UserResponse{
		ID:        value.ID.String(),
		Username:  value.Username,
		FullName:  value.FullName,
		IsActive:  value.IsActive,
		Role:      RoleToChoice(value.Role), // convert enum to ChoiceBase
		CarPark:   carPark,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
}

func ToUserResponseList(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	return responses
}

// UserPaginatedResponse is a struct for paginated API responses
type UserPaginatedResponse struct {
	Rows  []UserResponse `json:"rows" validate:"required"`  // generic slice of any type
	Page  int            `json:"page" validate:"required"`  // current page
	Limit int            `json:"limit" validate:"required"` // items per page
} // @name UserPaginatedResponse

type UserMessageResponse struct {
	Message string        `json:"message"`
	Data    *UserResponse `json:"data" nullable:"true" extensions:"x-nullable=true"`
} // @name UserMessageResponse
