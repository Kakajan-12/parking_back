package repository

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"backend/contrib/models"
)

// CameraRepoFilter represents optional filters for listing cameras
type CameraRepoFilter struct {
	Search         *string // general search across multiple columns
	Type           *models.CameraType
	IncludeDeleted *bool // include soft-deleted records
	ExcludeID      *int  // exclude specific camera ID
	Limit          *int
	Offset         *int
	OrderBy        *string
}

type CameraRepository struct {
	db *gorm.DB
}

// NewCameraRepository creates a new CameraRepository
func NewCameraRepository(db *gorm.DB) *CameraRepository {
	return &CameraRepository{db: db}
}

// buildCameraQuery applies filters from CameraRepoFilter
func (r *CameraRepository) buildCameraQuery(filter *CameraRepoFilter) *gorm.DB {
	query := r.db.Model(&models.Camera{})

	if filter == nil {
		return query
	}

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	if filter.ExcludeID != nil {
		query = query.Where("id != ?", *filter.ExcludeID)
	}

	if filter.IncludeDeleted != nil && *filter.IncludeDeleted {
		query = query.Unscoped()
	}

	if filter.Search != nil && *filter.Search != "" {
		likePattern := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", likePattern)
	}

	// Order by
	if filter.OrderBy != nil && *filter.OrderBy != "" {
		orderDir := "ASC"
		column := *filter.OrderBy
		if strings.HasPrefix(column, "-") {
			orderDir = "DESC"
			column = column[1:]
		}
		validColumns := map[string]bool{
			"id":         true,
			"name":       true,
			"type":       true,
			"created_at": true,
			"updated_at": true,
		}
		if validColumns[column] {
			query = query.Order(fmt.Sprintf("%s %s", column, orderDir))
		}
	} else {
		query = query.Order("created_at DESC")
	}

	return query
}

// CameraCreate creates a new camera
func (r *CameraRepository) CameraCreate(name string, cameraType models.CameraType, channelName *string, channelToken *string) (*models.Camera, error) {
	camera := &models.Camera{
		Name:         name,
		Type:         cameraType,
		ChannelName:  channelName,
		ChannelToken: channelToken,
	}

	if err := r.db.Create(camera).Error; err != nil {
		return nil, err
	}
	return camera, nil
}

// CameraGetByID retrieves a camera by ID
func (r *CameraRepository) CameraGetByID(id int64, includeDeleted bool) (*models.Camera, error) {
	var camera models.Camera
	query := r.db.Model(&models.Camera{}).Where("id = ?", id)
	if includeDeleted {
		query = query.Unscoped()
	}
	if err := query.First(&camera).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &camera, nil
}

// CameraList retrieves cameras based on filters
func (r *CameraRepository) CameraList(filter *CameraRepoFilter) ([]models.Camera, error) {
	var cameras []models.Camera
	query := r.buildCameraQuery(filter)

	if filter != nil {
		if filter.Limit != nil && *filter.Limit > 0 {
			query = query.Limit(*filter.Limit)
		}
		if filter.Offset != nil && *filter.Offset >= 0 {
			query = query.Offset(*filter.Offset)
		}
	}

	if err := query.Find(&cameras).Error; err != nil {
		return nil, err
	}
	return cameras, nil
}

// CameraCount returns the number of cameras matching the filter
func (r *CameraRepository) CameraCount(filter *CameraRepoFilter) (int64, error) {
	var count int64
	query := r.buildCameraQuery(filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CameraExist checks if a camera exists matching the filter
func (r *CameraRepository) CameraExist(filter *CameraRepoFilter) (bool, error) {
	query := r.buildCameraQuery(filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CameraUpdate updates an existing camera
func (r *CameraRepository) CameraUpdate(camera *models.Camera) error {
	return r.db.Model(&models.Camera{}).Where("id = ?", camera.ID).Updates(camera).Error
}

// CameraDelete soft-deletes a camera by ID
func (r *CameraRepository) CameraDelete(id int64) error {
	return r.db.Delete(&models.Camera{}, id).Error
}
