package repo

import (
	"database-example/model"

	"gorm.io/gorm"
)

type TokenVerificatonRepository struct {
	DatabaseConnection *gorm.DB
}

func NewTokenVerificatinRepository(databaseConnection *gorm.DB) *TokenVerificatonRepository {
	return &TokenVerificatonRepository{DatabaseConnection: databaseConnection}
}

func (repo *TokenVerificatonRepository) FindById(id string) (model.VerificationToken, error) {
	user := model.VerificationToken{}
	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)
	if dbResult != nil {
		return user, dbResult.Error
	}
	return user, nil
}

func (repo *TokenVerificatonRepository) CreateVerificatonToken(user *model.VerificationToken) error {
	dbResult := repo.DatabaseConnection.Create(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
