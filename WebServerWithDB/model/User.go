package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name" gorm:"not null;type:string"`
	Major string    `json:"major"`
}

func (user *User) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}
