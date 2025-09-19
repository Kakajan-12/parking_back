package schema

import (
	"backend/contrib/decimal"
	"backend/contrib/models"
	"time"
)

type TariffCreateInput struct {
	Name        string          `json:"name" validate:"required,max=255"`
	Duration    int             `json:"duration" validate:"required,gt=0"`
	IsActive    bool            `json:"is_active"`
	PriceAmount decimal.Decimal `json:"price_amount" validate:"required,decimal_gt0"`
} // @name TariffCreateInput

// TariffUpdateInput defines the input data for updating a tariff
type TariffUpdateInput struct {
	Name        *string          `json:"name" validate:"omitempty,max=255"`
	Duration    *int             `json:"duration" validate:"omitempty,gt=0"`
	IsActive    *bool            `json:"is_active"`
	PriceAmount *decimal.Decimal `json:"price_amount" validate:"omitempty,decimal_gt0"`
} // @name TariffUpdateInput

// TariffResponse is used for API responses.
type TariffResponse struct {
	ID          int64           `json:"id" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Duration    int             `json:"duration" validate:"required"`
	IsActive    bool            `json:"isActive" validate:"required"`
	PriceAmount decimal.Decimal `json:"PriceAmount" validate:"required"`
	Currency    string          `json:"currency" validate:"required"`
	CreatedAt   time.Time       `json:"createdAt" validate:"required"`
	UpdatedAt   *time.Time      `json:"updatedAt"` // optional, can be nil
} // @name TariffResponse

// ToTariffResponse maps a single DB model to response model
func ToTariffResponse(t models.Tariff) TariffResponse {
	return TariffResponse{
		ID:          t.ID,
		Name:        t.Name,
		Duration:    t.Duration,
		IsActive:    t.IsActive,
		PriceAmount: t.PriceAmount,
		Currency:    t.Currency,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// ToTariffResponseList maps a slice of DB models to a slice of response models
func ToTariffResponseList(tariffs []models.Tariff) []TariffResponse {
	responses := make([]TariffResponse, len(tariffs))
	for i, t := range tariffs {
		responses[i] = ToTariffResponse(t)
	}
	return responses
}

// TariffPaginatedResponse is a paginated response for tariffs
type TariffPaginatedResponse struct {
	Rows  []TariffResponse `json:"rows" validate:"required"`
	Page  int              `json:"page" validate:"required"`
	Limit int              `json:"limit" validate:"required"`
} // @name TariffPaginatedResponse

// TariffMessageResponse is used for responses that include a tariff and a message
type TariffMessageResponse struct {
	Message string          `json:"message"`
	Data    *TariffResponse `json:"data" nullable:"true"`
} // @name TariffMessageResponse
