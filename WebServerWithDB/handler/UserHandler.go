package handler

import (
	"database-example/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	encounterService := service.NewUserService(db)

	return &UserHandler{
		UserService: encounterService,
	}
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {

}
