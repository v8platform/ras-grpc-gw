package domain

type User struct {
	ID           int32
	UUID         string
	Username     string
	PasswordHash string
	Email        string
	IsAdmin      bool
}
