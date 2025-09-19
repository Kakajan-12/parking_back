package repository

import (
	"backend/contrib/models"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepoFilter struct {
	Search           *string          // general search across multiple columns
	Username         *string          // specific filter
	FullName         *string          // specific filter
	Role             *models.RoleType // specific filter
	IsActive         *bool            // specific filter
	IncludeSuperuser *bool            // optional: include superusers
	IncludeDeleted   *bool            // new property
	Limit            *int             // optional pagination
	Offset           *int             // optional pagination
	ExcludeID        *uuid.UUID
	OrderBy          *string
}

// buildUserQuery applies filters from UserFilter to a gorm.DB query
func (r *UserRepository) buildUserQuery(filter *UserRepoFilter) *gorm.DB {
	query := r.db.Model(&models.User{})

	if filter == nil {
		return query
	}

	// Specific filters
	if filter.Username != nil && *filter.Username != "" {
		query = query.Where("username = ?", *filter.Username)
	}

	if filter.FullName != nil && *filter.FullName != "" {
		query = query.Where("full_name = ?", *filter.FullName)
	}

	if filter.Role != nil && *filter.Role != "" {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Filter superusers only if IncludeSuperuser is nil or false
	if filter.IncludeSuperuser == nil || !*filter.IncludeSuperuser {
		query = query.Where("is_superuser = false")
	}
	// Filter deleted
	if filter.IncludeDeleted != nil && *filter.IncludeDeleted {
		// include soft-deleted records
		query = query.Unscoped()
	}

	if filter.ExcludeID != nil {
		query = query.Where("id != ?", *filter.ExcludeID)
	}

	// General search across multiple columns
	if filter.Search != nil && *filter.Search != "" {
		likePattern := "%" + *filter.Search + "%"
		query = query.Where(
			r.db.Where("username ILIKE ?", likePattern).
				Or("full_name ILIKE ?", likePattern),
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
			"username":   true,
			"full_name":  true,
			"created_at": true,
			"updated_at": true,
			"role":       true,
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

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// UserCreate Create a new user
func (r *UserRepository) UserCreate(
	username string,
	password string, // mandatory
	fullName string,
	role models.RoleType,
	isActive bool,
	carPark *models.CarParkType,
	isSuperuser *bool, // caller can pass, but default enforced
) (*models.User, error) {

	// Ensure password is not empty
	if password == "" {
		return nil, fmt.Errorf("password must be provided")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// default false if nil
	isSuperuserVal := false
	if isSuperuser != nil {
		isSuperuserVal = *isSuperuser
	}
	user := &models.User{
		ID:          uuid.New(),
		Username:    username,
		FullName:    fullName,
		Role:        role,
		IsActive:    isActive,
		CarPark:     carPark,
		IsSuperuser: isSuperuserVal,
		Password:    string(hashedPassword),
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UserGetByID Get user by ID
func (r *UserRepository) UserGetByID(id uuid.UUID, excludeDeleted bool) (*models.User, error) {
	var user models.User
	query := r.db.Model(&models.User{}).Where("id = ?", id)
	if !excludeDeleted {
		// include soft-deleted records
		query = query.Unscoped()
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UserList lists all users with optional filters and optional pagination.
// Filters map can include: "username", "role", "is_active", "limit", "offset"
// returns a list of users based on filters, search, and pagination
func (r *UserRepository) UserList(filter *UserRepoFilter) ([]models.User, error) {
	var users []models.User

	query := r.buildUserQuery(filter)

	// Optional pagination
	if filter.Limit != nil && *filter.Limit > 0 {
		query = query.Limit(*filter.Limit)
	}

	if filter.Offset != nil && *filter.Offset >= 0 {
		query = query.Offset(*filter.Offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UserCount retrieve users count by filter
func (r *UserRepository) UserCount(filter *UserRepoFilter) (int64, error) {
	var count int64
	query := r.buildUserQuery(filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UserExist checks if any user exists matching the given filter
func (r *UserRepository) UserExist(filter *UserRepoFilter) (bool, error) {
	query := r.buildUserQuery(filter) // reuse the same filter logic

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UserUpdate Update user
func (r *UserRepository) UserUpdate(u *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", u.ID).Updates(u).Error
}

// UserDelete Delete user by ID
func (r *UserRepository) UserDelete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UserAuthenticate checks username and password
func (r *UserRepository) UserAuthenticate(username string, password string) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("username = ?", username).
		Where("deleted_at IS NULL"). // ignore soft-deleted users
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // user not found
	} else if err != nil {
		return nil, err // DB error
	}

	// Compare password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil // password does not match
	}

	// Authentication successful
	return &user, nil
}
