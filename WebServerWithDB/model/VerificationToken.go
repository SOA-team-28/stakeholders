package model

import (
	"time"
)

type VerificationToken struct {
	userId            float64
	tokenCreationTime *time.Time
	TokenData         string
}
