package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"

	"gorm.io/gorm"
)

type VerificationTokenService struct {
	VerificationTokenRepo *repo.TokenVerificatonRepository
}

func NewVerificationTokenService(db *gorm.DB) *VerificationTokenService {
	repo := repo.NewTokenVerificatinRepository(db)
	return &VerificationTokenService{VerificationTokenRepo: repo}
}

func (service *VerificationTokenService) FindVerificationTokenByUser(id int) (*model.VerificationToken, error) {
	user, err := service.VerificationTokenRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &user, nil
}

func (service *VerificationTokenService) Create(user *model.VerificationToken) error {
	err := service.VerificationTokenRepo.CreateVerificatonToken(user)
	if err != nil {
		return err
	}
	return nil
}
