package repository

import (
	"backend/contrib/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

type UserSessionRepoFilter struct {
	UserID    *uuid.UUID
	Limit     *int // optional pagination
	Offset    *int // optional pagination
	ExcludeID *uuid.UUID
	OrderBy   *string
}

// buildUserQuery applies filters from UserFilter to a gorm.DB query
func (r *UserSessionRepository) buildUserSessionQuery(filter *UserSessionRepoFilter) *gorm.DB {
	query := r.db.Model(&models.UserSession{})

	if filter == nil {
		return query
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ExcludeID != nil {
		query = query.Where("id != ?", *filter.ExcludeID)
	}

	if filter.OrderBy != nil && *filter.OrderBy != "" {
		orderDir := "ASC"
		column := *filter.OrderBy

		if strings.HasPrefix(column, "-") {
			orderDir = "DESC"
			column = column[1:] // remove leading minus
		}

		// Optional: validate column name
		validColumns := map[string]bool{
			"id":         true,
			"expire_at":  true,
			"revoked_at": true,
			"created_at": true,
		}
		if validColumns[column] {
			query = query.Order(fmt.Sprintf("%s %s", column, orderDir))
		}
	} else {
		// Default ordering
		query = query.Order("created_at DESC")
	}

	return query
}

// UserSessionCreate Create new user session
func (r *UserSessionRepository) UserSessionCreate(
	userID uuid.UUID,
	expireAt time.Time,
	ipAddress *string,
	userAgent *string,
) (*models.UserSession, error) {
	// Create new UserSession instance
	userSession := &models.UserSession{
		ID:        uuid.New(), // generate new UUID
		UserID:    userID,
		IpAddress: ipAddress,
		UserAgent: userAgent,
		ExpireAt:  expireAt,
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := r.db.Create(userSession).Error; err != nil {
		return nil, err
	}

	return userSession, nil
}

// UserSessionGetByID Get user session by ID
func (r *UserSessionRepository) UserSessionGetByID(id uuid.UUID, withJoin bool) (*models.UserSession, error) {
	var userSession models.UserSession
	q := r.db.Model(&models.UserSession{})

	if withJoin {
		//q = q.Joins("INNER JOIN users ON users.id = user_sessions.user_id")
		q = q.Joins("User")
	}

	if err := q.First(&userSession, "user_sessions.id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userSession, nil
}

func (r *UserSessionRepository) UserSessionList(filter *UserSessionRepoFilter, withJoin bool) ([]models.UserSession, error) {
	var userSessions []models.UserSession
	query := r.buildUserSessionQuery(filter)
	if withJoin {
		query = query.Joins("User")
	}
	// Optional pagination
	if filter.Limit != nil && *filter.Limit > 0 {
		query = query.Limit(*filter.Limit)
	}

	if filter.Offset != nil && *filter.Offset >= 0 {
		query = query.Offset(*filter.Offset)
	}

	if err := query.Find(&userSessions).Error; err != nil {
		return nil, err
	}

	return userSessions, nil
}

// UserSessionUpdate Update user session
func (r *UserSessionRepository) UserSessionUpdate(u *models.UserSession) error {
	return r.db.Model(&models.UserSession{}).Where("id = ?", u.ID).Updates(u).Error
}

// UserSessionRevoke a session by setting RevokedAt
func (r *UserSessionRepository) UserSessionRevoke(id uuid.UUID) (int64, error) {
	now := time.Now()
	result := r.db.Model(&models.UserSession{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", now)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UserSessionRevokeAllByUserID revokes all sessions for a specific user
func (r *UserSessionRepository) UserSessionRevokeAllByUserID(userID uuid.UUID, exceptSessionID *uuid.UUID) (int64, error) {
	now := time.Now()

	// Update RevokedAt for all sessions matching the userID
	query := r.db.Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID)

	if exceptSessionID != nil {
		query = query.Where("id != ?", *exceptSessionID)
	}

	result := query.Update("revoked_at", &now)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
