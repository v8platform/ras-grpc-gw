package domain

type User struct {
	UUID         string
	Username     string
	PasswordHash string
	Email        string
	IsAdmin      bool
	Applications []string
}
