package repository

import (
	"backend/contrib/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// CarSessionRepository handles DB operations
type CarSessionRepository struct {
	db *gorm.DB
}

type CarSessionRepoFilter struct {
	Search    *string // general search across multiple columns
	CarNumber *string
	Limit     *int // optional pagination
	Offset    *int // optional pagination
	ExcludeID *int64
	OrderBy   *string
}

func (r *CarSessionRepository) buildCarSessionQuery(filter *CarSessionRepoFilter) *gorm.DB {
	query := r.db.Model(&models.CarSession{})

	if filter == nil {
		return query
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
			"created_at": true,
			"updated_at": true,
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

func NewCarSessionRepository(db *gorm.DB) *CarSessionRepository {
	return &CarSessionRepository{db: db}
}

func (r *CarSessionRepository) CarSessionCreate(dbObj *models.CarSession) (*models.CarSession, error) {
	if err := r.db.Create(dbObj).Error; err != nil {
		return nil, err
	}
	return dbObj, nil
}

func (r *CarSessionRepository) CarSessionGetByID(id int64) (*models.CarSession, error) {
	var dbObj models.CarSession
	if err := r.db.First(&dbObj, id).Error; err != nil {
		return nil, err
	}
	return &dbObj, nil
}

func (r *CarSessionRepository) CarSessionList(filter *CarSessionRepoFilter) ([]models.CarSession, error) {
	var rows []models.CarSession
	query := r.buildCarSessionQuery(filter)

	// Optional pagination
	if filter.Limit != nil && *filter.Limit > 0 {
		query = query.Limit(*filter.Limit)
	}

	if filter.Offset != nil && *filter.Offset >= 0 {
		query = query.Offset(*filter.Offset)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *CarSessionRepository) CarSessionCount(filter *CarSessionRepoFilter) (int64, error) {
	var count int64
	query := r.buildCarSessionQuery(filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CarSessionRepository) CarSessionUpdate(u *models.CarSession) error {
	return r.db.Model(&models.CarSession{}).Where("id = ?", u.ID).Updates(u).Error
}

func (r *CarSessionRepository) CarSessionDelete(id int) error {
	return r.db.Delete(&models.CarSession{}, id).Error
}

func (r *CarSessionRepository) CarSessionExist(filter *CarSessionRepoFilter) (bool, error) {
	query := r.buildCarSessionQuery(filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
