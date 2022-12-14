package user

import "github.com/jmoiron/sqlx"

type UserRepository interface {
	Create(user *User) (*User, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *User) (*User, error) {
	return nil, nil
}
