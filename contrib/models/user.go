package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/contrib/decimal"
)

// RoleType defines a custom type for user roles.
type RoleType string // @name RoleType

const (
	AdminRole      RoleType = "admin"
	OperatorRole   RoleType = "operator"
	AccountantRole RoleType = "accountant"
)

// User represents a user in the system.
type User struct {
	ID          uuid.UUID    `gorm:"column:id;type:uuid;primaryKey;not null"`
	Username    string       `gorm:"column:username;uniqueIndex;size:100;not null;"`
	FullName    string       `gorm:"column:full_name;size:100;not null"`
	Password    string       `gorm:"column:password;size:255;not null;"`
	IsActive    bool         `gorm:"column:is_active;default:true;"`
	IsSuperuser bool         `gorm:"column:is_superuser;default:false;"`
	Role        RoleType     `gorm:"column:role;type:varchar(20);not null;"`
	CarPark     *CarParkType `gorm:"column:car_park;type:varchar(20);"`

	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;index;"`
} // @name User

func (User) TableName() string {
	return "users"
}

type UserSession struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;not null"`
	IpAddress *string   `gorm:"column:ip_address;type:inet"`
	UserAgent *string   `gorm:"column:user_agent;type:text"`

	ExpireAt  time.Time  `gorm:"column:expire_at;not null;index"`
	RevokedAt *time.Time `gorm:"column:revoked_at;index"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`

	UserID uuid.UUID `gorm:"column:user_id;not null;index"`
	User   User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

type OperatorSession struct {
	ID int64 `gorm:"primaryKey;autoIncrement;column:id"`

	CarPark  CarParkType      `gorm:"column:car_park;not null"`
	LoginAt  time.Time        `gorm:"column:login_at;not null"`
	LogoutAt *time.Time       `gorm:"column:logout_at"`
	Money    *decimal.Decimal `gorm:"column:money;type:decimal(20,8);default:0.0"`
	Note     string           `gorm:"column:note;type:text"`
	Currency string           `gorm:"column:currency;size:20;not null"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;not null;"` // automatically set on insert
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`           // automatically set on update

	UserID uuid.UUID `gorm:"column:user_id;"`
	User   User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (OperatorSession) TableName() string {
	return "operator_sessions"
}
