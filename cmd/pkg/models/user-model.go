package models

import (
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
	db, err := Connection()
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
	db, err := Connection()
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
func RequestAdminStatus(userID int) error {
	fmt.Print(userID)
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	var currentStatus sql.NullString
	query := "SELECT admin_request_status FROM users WHERE id = ?"
	err = db.QueryRow(query, userID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no user found with ID %d", userID)
		}
		return err
	}

	if currentStatus.String == "pending" {
		return fmt.Errorf("admin request already pending for user ID %d", userID)
	}

	// Prepare the SQL statement to update the status
	updateQuery := "UPDATE users SET admin_request_status = ? WHERE id = ?"
	result, err := db.Exec(updateQuery, "pending", userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}
func GetAdminRequests() ([]types.User, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, username, admin_request_status FROM users WHERE admin_request_status IS NOT NULL"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Username, &user.AdminRequestStatus); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func ApproveAdminRequest(userID string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Update the admin_request_status
	query := "UPDATE users SET admin_request_status = ?, isAdmin = ? WHERE id = ?"
	result, err := tx.Exec(query, "approved", "1", userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no user found with ID %v", userID)
	}

	return tx.Commit()
}

func DenyAdminRequest(userID string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Prepare the SQL statement
	query := "UPDATE users SET admin_request_status = ? WHERE id = ?"

	// Execute the SQL statement
	result, err := db.Exec(query, "denied", userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %v", userID)
	}

	return nil
}

func GetAllAdmins() ([]types.User, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, username, isAdmin FROM users where isAdmin = 1`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adminList []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Username, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		adminList = append(adminList, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adminList, nil
}

func RemoveFromAdmin(userID string) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	query := `UPDATE users SET isAdmin = 0 WHERE id = ?`
	_, err = db.Exec(query, userID)

	if err != nil {
		return err
	}

	return nil
}
