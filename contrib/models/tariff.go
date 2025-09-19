package models

import (
	"backend/contrib/decimal"
	"time"
)

type Tariff struct {
	ID          int64           `gorm:"primaryKey;autoIncrement;column:id"`
	Name        string          `gorm:"column:name;size:255;not null"`
	Duration    int             `gorm:"column:duration;not null;uniqueIndex"`
	IsActive    bool            `gorm:"column:is_active;default:true;"`
	PriceAmount decimal.Decimal `json:"price"`
	Currency    string          `gorm:"column:currency;size:20;not null"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt   *time.Time      `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
}
