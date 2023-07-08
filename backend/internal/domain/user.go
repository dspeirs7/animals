package domain

import "context"

type User struct {
	Username string
	Password string
}

type UserRepository interface {
	GetUser(ctx context.Context, username string) (*User, error)
}
