package models

import (
	"bookstore/cmd/types"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GetBooks() ([]types.Book, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM books`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookList []types.Book
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Quantity, &book.Available, &book.CurrentBorrowedCount, &book.AvailableQuantity)
		if err != nil {
			return nil, err
		}
		bookList = append(bookList, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookList, nil
}

func GetBookById(id string) (types.Book, error) {
	db, err := Connection()
	if err != nil {
		return types.Book{}, err
	}

	query := `SELECT * FROM books WHERE ID = ?`
	row := db.QueryRow(query, id)

	var book types.Book
	err = row.Scan(&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Quantity, &book.Available, &book.CurrentBorrowedCount, &book.AvailableQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Book{}, errors.New("book not found")
		}
		return types.Book{}, err
	}

	return book, nil
}
func AddBook(book types.Book) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO books(isbn, title, author, quantity)
		VALUES(?,?,?,?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)
	`

	_, err = db.Exec(query, book.ISBN, book.Title, book.Author, book.Quantity)
	if err != nil {
		return err
	}

	return nil
}
func DeleteBook(id string) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	// Check if there are any borrowed copies of the book
	query := `SELECT COUNT(*) FROM checkouts WHERE book_id = ? AND return_date IS NULL`
	var borrowedCount int
	err = db.QueryRow(query, id).Scan(&borrowedCount)
	if err != nil {
		return err
	}

	if borrowedCount > 0 {
		return errors.New("book cannot be deleted because there are borrowed copies")
	}

	query = `DELETE FROM books WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateBook(id string, book types.Book) error {
	db, err := Connection()
	if err != nil {
		return err
	}

	// Check the number of borrowed copies of the book
	query := `SELECT COUNT(*) FROM checkouts WHERE book_id = ? AND return_date IS NULL`
	var borrowedCount int
	err = db.QueryRow(query, id).Scan(&borrowedCount)
	if err != nil {
		return err
	}
	var bookQty, _ = strconv.Atoi(book.Quantity)
	// Ensure the new quantity is not less than the number of borrowed copies
	if bookQty < borrowedCount {
		return errors.New("quantity cannot be less than the number of borrowed copies")
	}

	query = `
		UPDATE books
		SET quantity = ?
		WHERE id = ?
	`
	_, err = db.Exec(query, book.Quantity, id)
	if err != nil {
		return err
	}

	return nil
}

func SearchBooks(query string) ([]types.Book, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Adjust the query to search for books matching the query string
	searchQuery := "%" + query + "%"
	sqlStatement := `SELECT * FROM books WHERE (title LIKE ? OR author LIKE ?) AND available = true`
	rows, err := db.Query(sqlStatement, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Quantity, &book.Available, &book.CurrentBorrowedCount, &book.AvailableQuantity)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
func CreateCheckoutRequest(checkout types.CheckoutRequest) error {
	db, err := Connection()
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
	db, err := Connection()
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
	db, err := Connection()
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
	query = `UPDATE books SET available_quantity = available_quantity + 1 WHERE id = ?`
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

	db, err := Connection()
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
	db, err := Connection()
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
	_, err = tx.Exec("UPDATE books SET available_quantity = available_quantity - 1 WHERE id = (SELECT book_id FROM checkouts WHERE id = ?)", checkoutID)
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
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE checkouts SET checkout_status = 'denied' WHERE id = ?", checkoutID)
	if err != nil {
		return err
	}

	return nil
}
