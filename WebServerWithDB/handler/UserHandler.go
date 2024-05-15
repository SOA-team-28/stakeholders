package handler

import (
	"database-example/repo"
	"database-example/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository) *UserHandler {
	userService := service.NewUserService(db, tokenRepo)
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {

}
