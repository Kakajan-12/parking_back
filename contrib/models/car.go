package models

import (
	"time"

	"backend/contrib/decimal"
)

// CarParkType defines a custom type for car parks.
type CarParkType string // @name CarParkType

const (
	Park1 CarParkType = "P1"
	Park2 CarParkType = "P2"
	Park3 CarParkType = "P3"
	Park4 CarParkType = "P4"
)

type CarSessionEventType string // @name CarSessionEventType

const (
	EntryEvent CarSessionEventType = "entry"
	ExitEvent  CarSessionEventType = "exit"
)

// CarSessionStatusType defines a custom type for car session.
type CarSessionStatusType string // @name CarSessionStatusType

// Car represents the database model
type Car struct {
	ID        int64   `gorm:"primaryKey;autoIncrement"`
	CarNumber string  `gorm:"column:car_number;size:255;not null;uniqueIndex"`
	OwnerName *string `gorm:"column:owner_name;size:255;"`
	IsStaff   bool    `gorm:"column:is_staff;default:false;"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
}

func (Car) TableName() string {
	return "cars"
}

type CarSubscription struct {
	ID          int64           `gorm:"primaryKey;autoIncrement"`
	TotalAmount decimal.Decimal `gorm:"column:total_payment_amount;type:decimal(20,8);default:0.0"`
	Currency    string          `gorm:"column:currency;size:20;not null"`
	IsPaid      bool            `gorm:"column:is_paid;not null;default:false"`
	StartTime   *time.Time      `gorm:"column:start_time"`
	EndTime     *time.Time      `gorm:"column:end_time"`
	IsActive    bool            `gorm:"column:is_active;default:true;"`
	RevokedAt   *time.Time      `gorm:"column:revoked_at;index"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt   *time.Time      `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
	CarID       int64           `gorm:"column:car_id;not null;"`
	Car         Camera          `gorm:"foreignKey:CarID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (CarSubscription) TableName() string {
	return "car_subscriptions"
}

// CarSession represents the database model
type CarSession struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	Currency    string  `gorm:"column:currency;size:20;not null"`
	Status      string  `gorm:"column:status;size:100"`
	Reason      string  `gorm:"column:reason"`
	IsPaid      bool    `gorm:"column:is_paid;not null;default:false"`
	CarID       int64   `gorm:"column:car_id;not null;"`
	Car         Car     `gorm:"foreignKey:CarID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	TotalAmount *decimal.Decimal `gorm:"column:total_payment_amount;type:decimal(20,8);default:0.0"`
	StartTime   *time.Time       `gorm:"column:start_time"`
	EndTime     *time.Time       `gorm:"column:end_time"`
	CarPark     *CarParkType     `gorm:"column:car_park;size:100"`
	Duration    *decimal.Decimal `gorm:"column:duration;type:decimal(20,8);default:0.0"`

	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt     *time.Time `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
	InvalidatedAt *time.Time `gorm:"column:invalidated_at"`
}

func (CarSession) TableName() string {
	return "car_sessions"
}

type CarSessionEvent struct {
	ID        int64                  `gorm:"primaryKey;autoIncrement"`
	EventType CarSessionEventType    `gorm:"column:event_type;not null;size-20"`
	ExtraData map[string]interface{} `gorm:"column:extra_data;type:jsonb"`
	ImageUrl  *string                `gorm:"column:image_url"`

	CarSessionID int64      `gorm:"column:car_session_id;primaryKey;autoIncrement"`
	CarSession   CarSession `gorm:"foreignKey:CarSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CameraToken *string `gorm:"column:camera_token;size:255"`
	CameraID int64  `gorm:"column:camera_id;not null;"`
	Camera   Camera `gorm:"foreignKey:CameraID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
}
