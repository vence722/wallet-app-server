package repository

import (
	"database/sql"
	"time"
	"wallet-app-server/app/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User repository interface
type IUserRepository interface {
	GetUserByID(db *gorm.DB, userID string) (entity.User, error)
	GetUserByName(db *gorm.DB, userName string) (entity.User, error)
	CreateUserActivity(db *gorm.DB, userActType string, userActDetail string, userWalletID string, userActTime time.Time) error
}

// User repository instance
var UserRepository IUserRepository = &userRepositoryImpl{}

// User repository implementation
type userRepositoryImpl struct{}

// Get user by user_id
// If not found, return gorm.ErrRecordNotFound
func (ur *userRepositoryImpl) GetUserByID(db *gorm.DB, userID string) (entity.User, error) {
	var user entity.User
	err := db.Where("user_id = ?", userID).Find(&user).Error
	return user, err
}

// Get user by user_name
// If not found, return gorm.ErrRecordNotFound
func (ur *userRepositoryImpl) GetUserByName(db *gorm.DB, userName string) (entity.User, error) {
	var user entity.User
	err := db.Where("user_name = ?", userName).Find(&user).Error
	return user, err
}

// Create user activity
func (ur *userRepositoryImpl) CreateUserActivity(db *gorm.DB, userActType string, userActDetail string, userWalletID string, userActTime time.Time) error {
	userActivity := entity.UserActivity{
		UserActID:     uuid.New().String(),
		UserActType:   userActType,
		UserActDetail: userActDetail,
		UserActTime:   userActTime,
	}
	if userWalletID != "" {
		userActivity.UserWalletID = sql.NullString{String: userWalletID, Valid: true}
	}
	return db.Create(&userActivity).Error
}
