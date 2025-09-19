package repository

import (
	"backend/contrib/models"

	"gorm.io/gorm"
)

type MacUserRepository struct {
	db *gorm.DB
}

func NewMacUserRepository(db *gorm.DB) *MacUserRepository {
	return &MacUserRepository{db: db}
}

// MacUserCreate Create MacUser
func (r *MacUserRepository) MacUserCreate(m *models.MacUser) error {
	return r.db.Create(m).Error
}

// MacUserGetByID Get MacUser by ID
func (r *MacUserRepository) MacUserGetByID(id int64) (*models.MacUser, error) {
	var macUser models.MacUser
	if err := r.db.First(&macUser, id).Error; err != nil {
		return nil, err
	}
	return &macUser, nil
}

// MacUserList List all MacUsers
func (r *MacUserRepository) MacUserList() ([]models.MacUser, error) {
	var macUsers []models.MacUser
	if err := r.db.Find(&macUsers).Error; err != nil {
		return nil, err
	}
	return macUsers, nil
}

// MacUserExist checks if a macUser exists by username
func (r *MacUserRepository) MacUserExist(id int) (bool, error) {
	var count int64
	if err := r.db.Model(&models.MacUser{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// MacUserUpdate Update MacUser
func (r *MacUserRepository) MacUserUpdate(m *models.MacUser) error {
	return r.db.Save(m).Error
}

// MacUserDelete Delete MacUser by ID
func (r *MacUserRepository) MacUserDelete(id int) error {
	return r.db.Delete(&models.MacUser{}, id).Error
}
