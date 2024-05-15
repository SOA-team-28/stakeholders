package model

import (
	"time"
)

type VerificationToken struct {
	UserId            int
	TokenCreationTime time.Time
	TokenData         string
}
