package schema

import (
	"backend/core"
	"time"

	"backend/contrib/models"
)

var cameraTypeLabels = map[models.CameraType]string{
	models.InsideCamera:  "Inside Camera",
	models.OutsideCamera: "Outside Camera",
}

func CameraTypeToChoice(value models.CameraType) core.ChoiceBase {
	return core.ChoiceBase{
		Value: string(value),
		Label: cameraTypeLabels[value],
	}
}

// CameraCreateInput defines the input data for creating a camera
type CameraCreateInput struct {
	Name         string            `json:"name" validate:"required,max=255"`
	Type         models.CameraType `json:"type" validate:"required,cameratype"`
	ChannelName  *string           `json:"channelName" validate:"max=255" nullable:"true"`
	ChannelToken *string           `json:"channelToken" validate:"max=255" nullable:"true"`
} // @name CameraCreateInput

// CameraUpdateInput defines the input data for updating a camera
type CameraUpdateInput struct {
	Name         *string            `json:"name,omitempty" validate:"omitempty,max=255"`
	Type         *models.CameraType `json:"type,omitempty" validate:"omitempty"`
	ChannelName  *string            `json:"channelName" validate:"omitempty,max=255" nullable:"true"`
	ChannelToken *string            `json:"channelToken" validate:"omitempty,max=255" nullable:"true"`
} // @name CameraUpdateInput

// CameraResponse represents a camera returned in API responses
type CameraResponse struct {
	ID           int64           `json:"id" validate:"required"`
	Name         string          `json:"name" validate:"required"`
	Type         core.ChoiceBase `json:"type" validate:"required"`
	ChannelName  *string         `json:"channelName;"  nullable:"true"`
	ChannelToken *string         `json:"channelToken;"  nullable:"true"`
	CreatedAt    time.Time       `json:"createdAt" validate:"required"`
	UpdatedAt    *time.Time      `json:"updatedAt,omitempty"` // optional
} // @name CameraResponse

// ToCameraResponse converts a models.Camera to CameraResponse
func ToCameraResponse(c models.Camera) CameraResponse {
	return CameraResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      CameraTypeToChoice(c.Type),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToCameraResponseList converts a slice of models.Camera to a slice of CameraResponse
func ToCameraResponseList(cameras []models.Camera) []CameraResponse {
	res := make([]CameraResponse, len(cameras))
	for i, c := range cameras {
		res[i] = ToCameraResponse(c)
	}
	return res
}

// CameraPaginatedResponse is a paginated response for cameras
type CameraPaginatedResponse struct {
	Rows  []CameraResponse `json:"rows" validate:"required"`
	Page  int              `json:"page" validate:"required"`
	Limit int              `json:"limit" validate:"required"`
} // @name CameraPaginatedResponse

// CameraMessageResponse is used for responses that include a camera and a message
type CameraMessageResponse struct {
	Message string          `json:"message"`
	Data    *CameraResponse `json:"data,omitempty"`
} // @name CameraMessageResponse
