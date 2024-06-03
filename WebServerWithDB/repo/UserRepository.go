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

func (repo *UserRepository) FindById(id int) (model.User, error) {
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
func (repo *UserRepository) Login(username string, password string) (string, error) {
	user := model.User{}
	fmt.Println("Username u repo :", username) // Debug print

	dbResult := repo.DatabaseConnection.First(&user, "username = ?", username)

	if dbResult.Error != nil {
		return "no username", dbResult.Error
	}

	fmt.Print("u repo")
	fmt.Print(user)
	if user.Password != password {
		return "no password", fmt.Errorf("Neispravna lozinka")
	}
	fmt.Print("u repo")
	fmt.Print(user)

	// Check if the user already has a token
	var existingToken model.VerificationToken
	existingToken, err := repo.TokenRepository.FindById(user.Id)

	// Decode and validate the existing token
	claims, err := decodeToken(existingToken.TokenData)
	if err == nil && claims["userId"] == float64(user.Id) && time.Now().Unix() < int64(claims["expiry"].(float64)) {
		// Token is valid, return the existing token
		return existingToken.TokenData, nil
	}

	// Token is either invalid or does not exist, generate a new token
	tokenString, err := generateToken(user.Id, user.Username)
	if err != nil {
		return "", fmt.Errorf("Greška prilikom generisanja tokena: %v", err)
	}
	fmt.Print("u repo generisan token")
	fmt.Print(tokenString)
	newToken := model.VerificationToken{
		UserId:            user.Id,
		TokenCreationTime: time.Now(),
		TokenData:         tokenString,
	}

	// Save the new token in the database
	if err := repo.TokenRepository.CreateVerificatonToken(&newToken); err != nil {
		return "", fmt.Errorf("Greška prilikom kreiranja verifikacionog tokena: %v", err)
	}

	return tokenString, nil
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

func decodeToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte("secreet")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("greška prilikom parsiranja tokena: %v", err)
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("nevažeći token")
	}
}
func (repo *UserRepository) FindByUsername(username string) (model.User, error) {
	user := model.User{}
	dbResult := repo.DatabaseConnection.First(&user, "username = ?", username)
	if dbResult != nil {
		return user, dbResult.Error
	}
	return user, nil
}
