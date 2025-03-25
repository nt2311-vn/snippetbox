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

	query := `insert into (name, email, hashed_password, created)
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
	return 0, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
