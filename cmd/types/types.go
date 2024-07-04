package types

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Book struct {
	ID        int    `json:"id"`
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Quantity  string `json:"quantity"`
	Available string `json:"available"`
}

type JsonResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	IsAdmin            string `json:"isAdmin"`
	AdminRequestStatus string `json:"admin_request_status"`
}

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	IsAdmin  string `json:"isAdmin"`

	jwt.StandardClaims
}

type CheckoutRequest struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	BookID         int       `json:"book_id"`
	CheckoutDate   time.Time `json:"checkout_date"`
	DueDate        time.Time `json:"due_date"`
	ReturnDate     time.Time `json:"return_date"`
	CheckoutStatus string    `json:"checkout_status"`
	Fine           float64   `json:"fine"`
}

type BorrowingHistoryElt struct {
	ID           int       `json:"id"`
	BookID       int       `json:"book_id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	CheckoutDate time.Time `json:"checkout_date"`
	DueDate      time.Time `json:"due_date"`
	ReturnDate   time.Time `json:"return_date"`
	Fine         float64   `json:"fine"`
}

type CheckoutRequestPageElement struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	BookID         int     `json:"book_id"`
	Username       string  `json:"username"`
	Title          string  `json:"title"`
	CheckoutStatus string  `json:"checkout_status"`
	Fine           float64 `json:"fine"`
}
