package services

import (
	"mvc/datamodels"
	"mvc/repositories"
)

type UserService interface {
	GetAll() []datamodels.User
	GetByID(id int) (datamodels.User, bool)
	GetByUsernameAndPassword(username, password string) (datamodels.User, bool)
	DeleteByID(id int64) bool

	Update(id int64, user datamodels.User) (datamodels.User, error)
	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	UpdateUsername(id int64, newUsername string) (datamodels.User, error)

	Create(userPassword string, user datamodels.User) (datamodels.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}
