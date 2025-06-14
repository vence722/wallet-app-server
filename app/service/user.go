package service

import (
	"net/http"
	"time"
	"wallet-app-server/app/config"
	"wallet-app-server/app/constant"
	"wallet-app-server/app/db"
	"wallet-app-server/app/logger"
	"wallet-app-server/app/redis"
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
		logger.Errorf("Failed to get user from DB, err: %s", err.Error())
		return "", http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Validate user password
	inputPassHash := util.HashPassword(password)
	logger.Debugf("inputPassHash: %s, user.UserHash: %s", inputPassHash, user.UserHash)
	if user.UserHash != inputPassHash {
		return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessagePasswordNotValid, nil)
	}
	// Generate access token
	accessToken := uuid.New().String()
	// Insert access token into Redis
	if err := redis.Client.Set(accessToken, user.UserID, time.Duration(config.Cfg.Server.SessionExpireTimeInSecs)*time.Second); err != nil {
		logger.Errorf("Failed to insert access token to Redis, err: %s", err.Error())
		return "", http.StatusBadRequest, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, nil)
	}
	// Create user activity
	repository.UserRepository.CreateUserActivity(db.DB, user.UserID, constant.UserActTypeLogin, "User login", "", time.Now())
	return accessToken, http.StatusOK, nil
}
