package handler

import (
	"database-example/service"

	"gorm.io/gorm"
)

type VerificationTokenHandler struct {
	VerificationTokenServiceService *service.VerificationTokenService
}

func NewVerificationTokenHandler(db *gorm.DB) *VerificationTokenHandler {
	verificationTokenService := service.NewVerificationTokenService(db)

	return &VerificationTokenHandler{
		VerificationTokenServiceService: verificationTokenService,
	}
}
