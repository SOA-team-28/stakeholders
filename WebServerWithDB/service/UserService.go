package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repo.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	repo := repo.NewUserRepository(db)
	return &UserService{UserRepo: repo}
}

func (service *UserService) FindUser(id string) (*model.User, error) {
	user, err := service.UserRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &user, nil
}

func (service *UserService) Create(user *model.User) error {
	err := service.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
