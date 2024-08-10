package models

import (
	"bookstore/cmd/pkg/utils"
	"database/sql"
	"fmt"
)

func RequestAdminStatus(userID int) error {
	fmt.Print(userID)
	db, err := utils.Connection()
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
