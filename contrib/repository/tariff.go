package repository

import (
	"backend/contrib/decimal"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"backend/contrib/models"
)

// TariffRepository handles DB operations for Tariff
type TariffRepository struct {
	db *gorm.DB
}

type TariffRepoFilter struct {
	Duration  *int
	Limit     *int // optional pagination
	Offset    *int // optional pagination
	ExcludeID *int64
	OrderBy   *string
}

// buildTariffQuery applies filters from UserFilter to a gorm.DB query
func (r *TariffRepository) buildTariffQuery(filter *TariffRepoFilter) *gorm.DB {
	query := r.db.Model(&models.Car{})

	if filter == nil {
		return query
	}
	// Specific filters
	if filter.Duration != nil {
		query = query.Where("duration = ?", *filter.Duration)
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

// NewTariffRepository creates a new Tariff repository
func NewTariffRepository(db *gorm.DB) *TariffRepository {
	return &TariffRepository{db: db}
}

// TariffCreate creates a new tariff
func (r *TariffRepository) TariffCreate(name string, duration int, isActive bool, priceAmount decimal.Decimal, currency string) (*models.Tariff, error) {
	dbObj := &models.Tariff{
		Name:        name,
		Duration:    duration,
		IsActive:    isActive,
		PriceAmount: priceAmount,
		Currency:    currency,
	}

	if err := r.db.Create(dbObj).Error; err != nil {
		return nil, err
	}
	return dbObj, nil
}

// TariffGetByID GetByID Get a tariff by ID
func (r *TariffRepository) TariffGetByID(id int64) (*models.Tariff, error) {
	var tariff models.Tariff
	if err := r.db.First(&tariff, id).Error; err != nil {
		return nil, err
	}
	return &tariff, nil
}

// TariffList List all cars with optional filters
func (r *TariffRepository) TariffList(filter *TariffRepoFilter) ([]models.Tariff, error) {
	var rows []models.Tariff
	query := r.buildTariffQuery(filter)

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

// TariffCount returns the number of tariffs matching the filter
func (r *TariffRepository) TariffCount(filter *TariffRepoFilter) (int64, error) {
	var count int64
	query := r.buildTariffQuery(filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// TariffExist checks if a camera exists matching the filter
func (r *TariffRepository) TariffExist(filter *TariffRepoFilter) (bool, error) {
	query := r.buildTariffQuery(filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// TariffUpdate Update tariff
func (r *TariffRepository) TariffUpdate(u *models.Tariff) error {
	return r.db.Model(&models.Tariff{}).Where("id = ?", u.ID).Updates(u).Error
}

// TariffDelete Delete a tariff by ID
func (r *TariffRepository) TariffDelete(id int) error {
	return r.db.Delete(&models.Tariff{}, id).Error
}
