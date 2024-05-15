package repo

import (
	"database-example/model"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
	TokenRepository    *TokenVerificatonRepository
}

func NewUserRepository(databaseConnection *gorm.DB, tokenRepo *TokenVerificatonRepository) *UserRepository {
	return &UserRepository{DatabaseConnection: databaseConnection, TokenRepository: tokenRepo}
}

func (repo *UserRepository) FindById(id string) (model.User, error) {
	user := model.User{}
	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)
	if dbResult != nil {
		return user, dbResult.Error
	}
	return user, nil
}

func (repo *UserRepository) CreateUser(user *model.User) error {
	dbResult := repo.DatabaseConnection.Create(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *UserRepository) Login(username, password string) (model.User, error) {
	user := model.User{}
	dbResult := repo.DatabaseConnection.Where("username = ?", username).First(&user)
	if dbResult.Error != nil {
		return user, dbResult.Error
	}

	if user.Password != password {
		return user, fmt.Errorf("Neispravna lozinka")
	}

	tokenString, err := generateToken(user.Id, user.Username)
	if err != nil {
		return user, fmt.Errorf("Greška prilikom generisanja tokena: %v", err)
	}

	token := model.VerificationToken{
		UserId:            user.Id,
		TokenCreationTime: time.Now(),
		TokenData:         tokenString,
	}

	if err := repo.TokenRepository.CreateVerificatonToken(&token); err != nil {
		return user, fmt.Errorf("Greška prilikom kreiranja verifikacionog tokena: %v", err)
	}

	return user, nil
}
func generateToken(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userID,
		"username": username,
		"expiry":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte("secreet")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("greška prilikom potpisivanja tokena: %v", err)
	}

	return tokenString, nil
}
