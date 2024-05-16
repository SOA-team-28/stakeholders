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

func NewUserService(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository) *UserService {
	userRepo := repo.NewUserRepository(db, tokenRepo)
	return &UserService{UserRepo: userRepo}
}

func (service *UserService) FindUser(id int) (*model.User, error) {
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

func (service *UserService) Login(username string, password string) (string, error) {
	token, err := service.UserRepo.Login(username, password)
	fmt.Print("u servisu")
	fmt.Print(token)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot login "))
	}
	return token, nil
}
