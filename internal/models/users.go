package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 18)
	if err != nil {
		return err
	}

	query := `insert into users (name, email, hashed_password, created)
	values (?, ?, ?, utc_timestamp())
	`
	_, err = um.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLErr *mysql.MySQLError

		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPass []byte

	query := `select id, hashed_password from users where email = ?`

	err := um.DB.QueryRow(query, email).Scan(&id, &hashedPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
