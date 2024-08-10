package models

import (
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"database/sql"
	"errors"
	"strconv"
)

func GetBooks(UserID int) ([]types.BookPageElement, error) {
	db, err := utils.Connection()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM books`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookList []types.BookPageElement
	for rows.Next() {

		var book types.BookPageElement
		err := rows.Scan(&book.Book.ID, &book.Book.ISBN, &book.Book.Title, &book.Book.Author, &book.Book.Quantity, &book.Book.Available, &book.Book.CurrentBorrowedCount, &book.Book.AvailableQuantity)
		if err != nil {
			return nil, err
		}
		var borrowedQuery = `SELECT COUNT(*) FROM checkouts WHERE book_id = ? AND return_date IS NULL AND user_id = ?`
		err = db.QueryRow(borrowedQuery, book.Book.ID, UserID).Scan(&book.Status)
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
	db, err := utils.Connection()
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
	db, err := utils.Connection()
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
	db, err := utils.Connection()
	if err != nil {
		return err
	}

	// Check if there are any borrowed copies of the book
	query := `SELECT COUNT(*) FROM checkouts WHERE book_id = ? AND return_date IS NULL AND checkout_status = 'approved'`
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
	db, err := utils.Connection()
	if err != nil {
		return err
	}

	// Check the number of borrowed copies of the book
	query := `SELECT current_borrowed_count FROM books WHERE id = ?`
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
	db, err := utils.Connection()
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
