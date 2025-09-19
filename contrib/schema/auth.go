package schema

import (
	"backend/contrib/models"
	"time"

	"github.com/google/uuid"
)

// LoginInput user login data
type LoginInput struct {
	Username string              `json:"username" validate:"required,max=150"`
	Password string              `json:"password" validate:"required,max=150"`
	CarPark  *models.CarParkType `json:"carPark"  validate:"omitempty,carparktype" nullable:"true"`
} // @name LoginInput

// LoginResponse defines the schema for a successful login
type LoginResponse struct {
	AccessToken string          `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." validate:"required"`
	Role        models.RoleType `json:"role" example:"operator" validate:"required"`
} // @name LoginResponse

type UserSessionResponse struct {
	ID        string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required" description:"Unique identifier of the session"`
	ExpireAt  time.Time  `json:"expireAt" description:"Session expiration timestamp" validate:"required"`
	RevokedAt *time.Time `json:"revokedAt,omitempty" description:"Time when the session was revoked"`
	IpAddress *string    `json:"ipAddress" example:"192.168.1.1" description:"IP address of the user during the session"`
	UserAgent *string    `json:"userAgent" example:"Mozilla/5.0" description:"User agent string of the client"`
	CreatedAt time.Time  `json:"createdAt" description:"Session creation timestamp" validate:"required"`
	UserID    uuid.UUID  `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required" description:"Associated user ID"`
} // @name UserSessionResponse

func ToUserSessionResponse(s models.UserSession) UserSessionResponse {
	return UserSessionResponse{
		ID:        s.ID.String(),
		ExpireAt:  s.ExpireAt,
		RevokedAt: s.RevokedAt,
		IpAddress: s.IpAddress,
		UserAgent: s.UserAgent,
		CreatedAt: s.CreatedAt,
		UserID:    s.UserID,
	}
}

func ToUserSessionResponseList(users []models.UserSession) []UserSessionResponse {
	responses := make([]UserSessionResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserSessionResponse(u)
	}
	return responses
}

type UserSessionExtendedResponse struct {
	ID          string       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required" description:"Unique identifier of the session"`
	ExpireAt    time.Time    `json:"expireAt" description:"Session expiration timestamp" validate:"required"`
	RevokedAt   *time.Time   `json:"revokedAt,omitempty" description:"Time when the session was revoked"`
	IpAddress   *string      `json:"ipAddress" example:"192.168.1.1" description:"IP address of the user during the session"`
	UserAgent   *string      `json:"userAgent" example:"Mozilla/5.0" description:"User agent string of the client"`
	CreatedAt   time.Time    `json:"createdAt" description:"Session creation timestamp" validate:"required"`
	MacUsername *string      `json:"macUsername" example:"mac-username" description:"Device MAC username"`
	MacPassword *string      `json:"macPassword" example:"mac-password" description:"Device MAC password"`
	User        UserResponse `json:"user" description:"Associated user"`
} // @name UserSessionExtendedResponse

func ToUserSessionExtendedResponse(s models.UserSession) UserSessionExtendedResponse {
	return UserSessionExtendedResponse{
		ID:        s.ID.String(),
		ExpireAt:  s.ExpireAt,
		RevokedAt: s.RevokedAt,
		IpAddress: s.IpAddress,
		UserAgent: s.UserAgent,
		CreatedAt: s.CreatedAt,
		User:      ToUserResponse(s.User),
	}
}

func ToUserSessionExtendedResponseList(users []models.UserSession) []UserSessionExtendedResponse {
	responses := make([]UserSessionExtendedResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserSessionExtendedResponse(u)
	}
	return responses
}

// UserSessionExtendedPaginatedResponse is a struct for paginated API responses
type UserSessionExtendedPaginatedResponse struct {
	Rows  []UserSessionExtendedResponse `json:"rows" validate:"required"`  // generic slice of any type
	Page  int                           `json:"page" validate:"required"`  // current page
	Limit int                           `json:"limit" validate:"required"` // items per page
} // @name UserSessionExtendedPaginatedResponse

type TokenPayload struct {
	UserID string    `json:"userId"`
	Role   string    `json:"role"`
	Jti    string    `json:"jti"`
	Iat    time.Time `json:"iat"`
	Exp    time.Time `json:"exp"`
	Iss    string    `json:"iss"`
	Aud    string    `json:"aud"`
}
