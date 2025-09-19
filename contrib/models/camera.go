package models

import (
	"time"

	"gorm.io/gorm"
)

// CameraType defines a custom type for camera types.
type CameraType string // @name CameraType

const (
	InsideCamera  CameraType = "inside"
	OutsideCamera CameraType = "outside"
)

// Camera represents a physical camera
type Camera struct {
	ID           int64           `gorm:"primaryKey;autoIncrement;column:id"`
	Name         string          `gorm:"not null;column:name;size:255"`
	Type         CameraType      `gorm:"not null;column:type;type:varchar(50)"`
	ChannelName  *string         `gorm:"column:channel_name;size:255"`
	ChannelToken *string         `gorm:"column:channel_token;size:255"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt    *time.Time      `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}
