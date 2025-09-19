package core

// ErrorResponse is a generic error payload
type ErrorResponse struct {
	Detail string `json:"detail" example:"Invalid request body"`
} // @name ErrorResponse

// BadRequestResponse represents a 401 Unauthorized error
type BadRequestResponse struct {
	Detail string `json:"detail" example:"BadRequest - invalid request"`
} // @name BadRequestResponse

// UnauthorizedResponse represents a 401 Unauthorized error
type UnauthorizedResponse struct {
	Detail string `json:"detail" example:"Unauthorized - Invalid token"`
} // @name UnauthorizedResponse

// PermissionDeniedResponse represents a 403 Permission denied error
type PermissionDeniedResponse struct {
	Detail string `json:"detail" example:"Permission denied"`
} // @name PermissionDeniedResponse

// NotFoundResponse represents a 404 Permission denied error
type NotFoundResponse struct {
	Detail string `json:"detail" example:"Not found"`
} // @name NotFoundResponse

// ValidationErrorResponse represents validation errors (422)
type ValidationErrorResponse struct {
	Detail string            `json:"detail" example:"Validation error"`
	Errors map[string]string `json:"errors,omitempty" example:"fieldName:this field is required"`
} // @name ValidationErrorResponse

// Implement the error interface
func (v *ValidationErrorResponse) Error() string {
	return v.Detail
}

// InternalServerErrorResponse represents a 500 Internal Server Error
type InternalServerErrorResponse struct {
	Detail string `json:"detail" example:"Internal Server Error"`
} // @name InternalServerErrorResponse

type CountResponse struct {
	Count int64 `json:"count" example:"12345" validate:"required"`
} // @name CountResponse

type MessagedResponse struct {
	Message string `json:"message"`
} // @name MessagedResponse

// ChoiceBase response type for enums
type ChoiceBase struct {
	Value string `json:"value" validate:"required"`
	Label string `json:"label" validate:"required"`
} // @name ChoiceBase

type CommonsModel struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
