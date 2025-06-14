package service

import (
	"net/http"
	"wallet-app-server/app/db"
	"wallet-app-server/app/repository"
	"wallet-app-server/app/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User service interface
type IUserService interface {
	Login(username string, password string) (string, int, error)
}

// User service instance
var UserService IUserService = &userServiceImpl{}

// User service implementation
type userServiceImpl struct{}

func (us *userServiceImpl) Login(username string, password string) (string, int, error) {
	// Try to fetch user record by user name in DB
	user, err := repository.UserRepository.GetUserByName(db.DB, username)
	if err != nil {
		// If user not found,
		if err == gorm.ErrRecordNotFound {
			return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageUserNotFound, nil)
		}
		// Other repository error
		return "", http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Validate user password
	if user.UserHash != util.HashPassword(password) {
		return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessagePasswordNotValid, nil)
	}
	// Generate session key
	sessionKey := uuid.New().String()
	// Insert session key into Redis
	return sessionKey, http.StatusOK, nil
}
