package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"backend/contrib/models"
)

// CarRepository handles DB operations
type CarRepository struct {
	db *gorm.DB
}

type CarRepoFilter struct {
	Search    *string // general search across multiple columns
	CarNumber *string
	Limit     *int // optional pagination
	Offset    *int // optional pagination
	ExcludeID *int64
	OrderBy   *string
}

// buildUserQuery applies filters from UserFilter to a gorm.DB query
func (r *CarRepository) buildCarQuery(filter *CarRepoFilter) *gorm.DB {
	query := r.db.Model(&models.Car{})

	if filter == nil {
		return query
	}
	// Specific filters
	if filter.CarNumber != nil && *filter.CarNumber != "" {
		query = query.Where("car_number = ?", *filter.CarNumber)
	}

	if filter.ExcludeID != nil {
		query = query.Where("id != ?", *filter.ExcludeID)
	}

	// General search across multiple columns
	if filter.Search != nil && *filter.Search != "" {
		likePattern := "%" + *filter.Search + "%"
		query = query.Where(
			r.db.Where("car_number ILIKE ?", likePattern),
		)
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

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{db: db}
}

// CarCreate creates a new car
func (r *CarRepository) CarCreate(carNumber string, isStaff bool, ownerName *string) (*models.Car, error) {
	dbObj := &models.Car{
		CarNumber: carNumber,
		IsStaff:   isStaff,
		OwnerName: ownerName,
	}

	if err := r.db.Create(dbObj).Error; err != nil {
		return nil, err
	}
	return dbObj, nil
}

// CarGetByID Get a car by ID
func (r *CarRepository) CarGetByID(id int64) (*models.Car, error) {
	var car models.Car
	if err := r.db.First(&car, id).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

// CarList List all cars with optional filters
func (r *CarRepository) CarList(filter *CarRepoFilter) ([]models.Car, error) {
	var rows []models.Car
	query := r.buildCarQuery(filter)

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

// CarCount returns the number of cars matching the filter
func (r *CarRepository) CarCount(filter *CarRepoFilter) (int64, error) {
	var count int64
	query := r.buildCarQuery(filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CarUpdate Update tariff
func (r *CarRepository) CarUpdate(u *models.Car) error {
	return r.db.Model(&models.Car{}).Where("id = ?", u.ID).Updates(u).Error
}

// CarDelete Delete a car entry
func (r *CarRepository) CarDelete(id int) error {
	return r.db.Delete(&models.Car{}, id).Error
}

// CarExist checks if a camera exists matching the filter
func (r *CarRepository) CarExist(filter *CarRepoFilter) (bool, error) {
	query := r.buildCarQuery(filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
