package service

import "net/http"

// User service interface
type IUserService interface {
	Login(username string, password string) (string, int, error)
}

// User service instance
var UserService IUserService = &userServiceImpl{}

// User service implementation
type userServiceImpl struct{}

func (us *userServiceImpl) Login(username string, password string) (string, int, error) {
	return "", http.StatusOK, nil
}
