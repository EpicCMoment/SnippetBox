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
	ID	int
	Name	string
	Email	string
	HashedPassword	[]byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}


func hashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil

}

func (um *UserModel) Insert(name, email, password string) error {

	hashedPassword, err := hashPassword(password)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = um.DB.Exec(stmt, name, email, hashedPassword)

	if err != nil {

		var mySQLError *mysql.MySQLError

		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}

		return err

	}

	return nil
}

func (um *UserModel) Exists(id int) error {

	panic("not implemented")	

}

func (um *UserModel) Authenticate(email, password string) (int, error) {

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?;`

	row := um.DB.QueryRow(stmt, email)

	var userId int
	var hashedPassword string

	err := row.Scan(&userId, &hashedPassword)

	if err != nil {
		return -1, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return -1, err
	}

	return userId, nil

}