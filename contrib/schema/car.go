package schema

import (
	"time"

	"backend/contrib/decimal"
	"backend/contrib/models"
	"backend/core"
)

var carParkLabels = map[models.CarParkType]string{
	models.Park3: "Parking Lot 3",
	models.Park4: "Parking Lot 4",
}

func CarParkToChoice(value string) core.ChoiceBase {
	label, ok := carParkLabels[models.CarParkType(value)]
	if !ok {
		// If unknown, just return the value itself
		label = value
	}

	return core.ChoiceBase{
		Value: value,
		Label: label,
	}
}

var carStatusLabels = map[models.CarSessionStatusType]string{}

func StatusToChoice(value string) core.ChoiceBase {
	label, ok := carStatusLabels[models.CarSessionStatusType(value)]
	if !ok {
		label = value
	}
	return core.ChoiceBase{
		Value: value,
		Label: label,
	}
}

type CarCreateInput struct {
	CarNumber string  `json:"carNumber" validate:"required,max=255"`
	OwnerName *string `json:"ownerName" nullable:"true" validate:"max=255"`
	IsStaff   bool    `json:"isStaff"`
} // @name CarCreateInput

// CarUpdateInput defines the input data for updating a camera
type CarUpdateInput struct {
	CarNumber *string `json:"carNumber,omitempty" validate:"omitempty,max=255"`
	OwnerName *string `json:"ownerName,omitempty" validate:"omitempty,max=255"`
	IsStaff   *bool   `json:"isStaff,omitempty" validate:"omitempty"`
} // @name CarUpdateInput

// CarResponse is used for API responses
type CarResponse struct {
	ID        int64      `json:"id" validate:"required"`
	CarNumber string     `json:"carNumber" validate:"required"`
	OwnerName *string    `json:"ownerName" nullable:"true"`
	IsStaff   bool       `json:"isStaff"   validate:"required"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt"` // optional, can be nil
} // @name CarResponse

// ToCarResponse Mapper function for CarResponse
func ToCarResponse(value models.Car) CarResponse {
	return CarResponse{
		ID:        value.ID,
		CarNumber: value.CarNumber,
		OwnerName: value.OwnerName,
		IsStaff:   value.IsStaff,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
}

func ToCarResponseList(cars []models.Car) []CarResponse {
	responses := make([]CarResponse, len(cars))
	for i, c := range cars {
		responses[i] = ToCarResponse(c)
	}
	return responses
}

// CarPaginatedResponse is a paginated response for cameras
type CarPaginatedResponse struct {
	Rows  []CarResponse `json:"rows" validate:"required"`
	Page  int           `json:"page" validate:"required"`
	Limit int           `json:"limit" validate:"required"`
} // @name CarPaginatedResponse

// CarMessageResponse is used for responses that include a camera and a message
type CarMessageResponse struct {
	Message string       `json:"message"`
	Data    *CarResponse `json:"data" nullable:"true"`
} // @name CarMessageResponse

type CarSessionResponse struct {
	ID          int64            `json:"id" validate:"required"`
	StartTime   *time.Time       `json:"startTime"`
	EndTime     *time.Time       `json:"endTime"`
	TotalAmount *decimal.Decimal `json:"totalAmount" nullable:"true"`
	Currency    string           `json:"currency" validate:"required"`
	Status      core.ChoiceBase  `json:"status" validate:"required"`
	Reason      string           `json:"reason"`
	ImageUrl    *string          `json:"imageUrl,omitempty"`
	CarPark     *core.ChoiceBase `json:"carPark" `
	Duration    *decimal.Decimal `json:"duration"`
	IsPaid      bool             `json:"isPaid" validate:"required"`
	CameraToken string           `json:"cameraToken"`
	CameraID    int64            `json:"cameraId"`
	CarID       int64            `json:"carId" validate:"required"`
	Car         *CarResponse     `json:"car" nullable:"true"`
	CreatedAt   time.Time        `json:"createdAt" validate:"required"`
	UpdatedAt   *time.Time       `json:"updatedAt"` // optional, can be nil
} // @name CarSessionResponse

type CarSessionPaginatedResponse struct {
	Rows  []CarSessionResponse `json:"rows" validate:"required"`
	Page  int                  `json:"page" validate:"required"`
	Limit int                  `json:"limit" validate:"required"`
} // @name CarSessionPaginatedResponse

type CarSessionMessageResponse struct {
	Message string              `json:"message"`
	Data    *CarSessionResponse `json:"data" nullable:"true"`
} // @name CarSessionMessageResponse

func ToCarSessionResponse(value models.CarSession) CarSessionResponse {
	var carPark *core.ChoiceBase
	if value.CarPark != nil {
		c := CarParkToChoice(string(*value.CarPark))
		carPark = &c
	}
	car := ToCarResponse(value.Car)
	return CarSessionResponse{
		ID:          value.ID,
		StartTime:   value.StartTime,
		EndTime:     value.EndTime,
		TotalAmount: value.TotalAmount,
		Status:      StatusToChoice(string(value.Status)),
		Reason:      value.Reason, 
		CarPark:     carPark,
		Duration:    value.Duration,
		IsPaid:      value.IsPaid, 
		CarID:       value.CarID,
		Car:         &car,
	}
}

func ToCarSessionResponseList(rows []models.CarSession) []CarSessionResponse {
	responses := make([]CarSessionResponse, len(rows))
	for i, c := range rows {
		responses[i] = ToCarSessionResponse(c)
	}
	return responses
}

type CapturedEventData struct {
	EventID          *string    `json:"EventId"`
	EventDescription *string    `json:"EventDescription"`
	EventComment     *string    `json:"EventComment"`
	ChannelName      *string    `json:"ChannelName"`
	CapturedTime     *time.Time `json:"captured_time"`
}


type CapturedEventDataE struct {
	EventID          *string    `json:"EventId"`
	EventDescription *string    `json:"EventDescription"`
	EventComment     *string    `json:"EventComment"`
	ChannelName      *string    `json:"ChannelName"`
	CapturedTime     *time.Time `json:"captured_time"`
	ChannelId        *string    `json:"ChannelId"`
}