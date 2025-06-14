package repository

import "wallet-app-server/app/entity"

type IUserRepository interface {
	GetUserByID(userID string) (entity.User, error)
}
