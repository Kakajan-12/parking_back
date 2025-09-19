package models

import (
	"time"

	"github.com/google/uuid"
)

type MacUser struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;column:id"`
	MacUsername string `gorm:"column:mac_username;not null"`
	MacPassword string `gorm:"column:mac_password;not null"`
	IsActive    bool   `gorm:"default:true;column:is_active"`
}

type MacUserAction struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id"`
	MacUserID int64     `gorm:"column:car_id;not null;"`
	MacUser   MacUser   `gorm:"foreignKey:MacUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UserID    uuid.UUID `gorm:"column:user_id;"`
	User      MacUser   `gorm:"foreignKey:MacUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
}

func (MacUser) TableName() string {
	return "mac_users"
}
