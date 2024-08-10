package models

import (
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateCheckoutRequest(checkout types.CheckoutRequest) error {
	db, err := utils.Connection()
	if err != nil {
		return err
	}
	defer db.Close()
	var checkoutDateStr string
	query := `
		SELECT checkout_date
		FROM checkouts
		WHERE user_id = ? AND book_id = ? AND return_date IS NULL
	`
	err = db.QueryRow(query, checkout.UserID, checkout.BookID).Scan(&checkoutDateStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if checkoutDateStr != "" {
		checkoutDate, err := time.Parse("2006-01-02 15:04:05", checkoutDateStr)
		if err != nil {
			return err
		}
		return fmt.Errorf("user has already borrowed this book on %s and hasn't returned it yet", checkoutDate)
	}

	// Insert checkout request into the database
	insertQuery := `
        INSERT INTO checkouts (user_id, book_id, checkout_date, due_date)
        VALUES (?, ?, ?, ?)
    `
	_, err = db.Exec(insertQuery, checkout.UserID, checkout.BookID, checkout.CheckoutDate, checkout.DueDate)
	if err != nil {
		return err
	}

	return nil
}
func GetCheckoutHistory(userID int) ([]types.BorrowingHistoryElt, error) {
	db, err := utils.Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
        SELECT checkouts.id, checkouts.book_id, books.title, books.author, checkouts.checkout_date, checkouts.due_date, checkouts.return_date, checkouts.fine
        FROM checkouts
        JOIN books ON checkouts.book_id = books.id
        WHERE checkouts.user_id = ? AND checkouts.checkout_status = 'approved'
    `

	rows, err := db.Query(query, userID)

	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var borrowingHistory []types.BorrowingHistoryElt
	for rows.Next() {
		var book types.BorrowingHistoryElt
		var checkoutDate, dueDate string
		var returnDate sql.NullString

		// Ensure the scan order matches the SELECT order
		err := rows.Scan(&book.ID, &book.BookID, &book.Title, &book.Author, &checkoutDate, &dueDate, &returnDate, &book.Fine)

		if err != nil {
			return nil, err
		}

		book.CheckoutDate, err = time.Parse("2006-01-02 15:04:05", checkoutDate)
		if err != nil {

			return nil, err
		}

		book.DueDate, err = time.Parse("2006-01-02 15:04:05", dueDate)
		if err != nil {

			return nil, err
		}

		if returnDate.Valid {
			book.ReturnDate, err = time.Parse("2006-01-02 15:04:05", returnDate.String)
			if err != nil {

				return nil, err
			}
		}

		// Explicitly set to nil if returnDate is not valid
		borrowingHistory = append(borrowingHistory, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return borrowingHistory, nil
}

func CheckinBook(checkoutID int) error {
	db, err := utils.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Update the return_date in checkouts table
	query := `UPDATE checkouts SET return_date = ? WHERE id = ?`
	_, err = tx.Exec(query, time.Now(), checkoutID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the book ID associated with the checkout
	var bookID int
	query = `SELECT book_id FROM checkouts WHERE id = ?`
	err = tx.QueryRow(query, checkoutID).Scan(&bookID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Increment the quantity in the books table
	query = `UPDATE books SET current_borrowed_count  = current_borrowed_count  + 1 WHERE id = ?`
	_, err = tx.Exec(query, bookID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func GetAllCheckouts() ([]types.CheckoutRequestPageElement, error) {

	db, err := utils.Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT c.id, u.username, b.title AS book_title, c.checkout_status 
		FROM checkouts c 
		INNER JOIN users u ON c.user_id = u.id 
		INNER JOIN books b ON c.book_id = b.id
	`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare a slice to hold the results
	var checkouts []types.CheckoutRequestPageElement

	// Iterate over the rows and populate the checkouts slice
	for rows.Next() {
		var checkout types.CheckoutRequestPageElement
		if err := rows.Scan(&checkout.ID, &checkout.Username, &checkout.Title, &checkout.CheckoutStatus); err != nil {
			return nil, err
		}
		checkouts = append(checkouts, checkout)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkouts, nil
}
func ApproveCheckout(checkoutID string) error {
	// Start a transaction
	db, err := utils.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the available_quantity is greater than 0
	var availableQuantity int
	err = tx.QueryRow(`
		SELECT b.available_quantity
		FROM books b
		JOIN checkouts c ON b.id = c.book_id
		WHERE c.id = ?`, checkoutID).Scan(&availableQuantity)
	if err != nil {
		return err
	}

	if availableQuantity <= 0 {
		return errors.New("cannot approve checkout: no available copies")
	}

	// Update the checkout status to 'approved'
	_, err = tx.Exec("UPDATE checkouts SET checkout_status = 'approved' WHERE id = ?", checkoutID)
	if err != nil {
		return err
	}

	// Decrement the available_quantity of the book in the books table
	_, err = tx.Exec("UPDATE books SET current_borrowed_count  = current_borrowed_count + 1 WHERE id = (SELECT book_id FROM checkouts WHERE id = ?)", checkoutID)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func DenyCheckout(checkoutID string) error {
	// Update the checkout status to 'denied'
	db, err := utils.Connection()
	if err != nil {
		return err
		// fmt.Print(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("UPDATE checkouts SET checkout_status = 'denied' WHERE id = ?", checkoutID)
	if err != nil {
		return err
		// fmt.Print(err.Error())
	}

	return nil
}
