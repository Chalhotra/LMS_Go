package models

import (
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RegisterUser(user types.User) error {
	db, err := utils.Connection()
	if err != nil {
		return err
	}

	username := user.Username
	passwd := user.Password
	password, err := hashPassword(passwd)
	if err != nil {
		return err
	}

	insertQuery := `INSERT INTO users(username, password) VALUES(?,?)`
	_, err = db.Exec(insertQuery, username, password)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return errors.New("username already exists")
		}
		return err
	}
	return nil
}

func GetUserByName(username string) (types.User, error) {
	var userInfo types.User
	db, err := utils.Connection()
	if err != nil {
		return userInfo, err
	}
	checkQuery := `SELECT id, username, password, isAdmin FROM users WHERE username = ?`
	row := db.QueryRow(checkQuery, username)

	err = row.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.IsAdmin)
	if err == sql.ErrNoRows {
		return userInfo, fmt.Errorf("user with the given username does not exist")
	}
	return userInfo, nil
}
