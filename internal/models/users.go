package models

import (
	"database/sql"
	"time"
)

type User struct {
	HashedPassword []byte
	Created        time.Time
	Name           string
	Email          string
	ID             int
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {
	return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
