package model

type User struct {
	Id                int
	Username          string
	Password          string
	Role              UserRole
	IsActive          bool
	Email             string
	VerificationToken string
	IsVerified        bool
}

type UserRole int

const (
	Administrator UserRole = iota
	Author
	Tourist
)
