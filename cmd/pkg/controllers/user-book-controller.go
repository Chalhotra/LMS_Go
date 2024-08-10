package controllers

import (
	"bookstore/cmd/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	book, err := models.GetBookById(params["id"])
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	user, ok := GetUserFromContext(r)
	if !ok {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=Unauthorized")
		return
	}
	bookList, err := models.GetBooks(user.ID)
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(bookList)
}
func SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Bad Request")
		return
	}

	// Fetch books based on the query
	books, err := models.SearchBooks(query)
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Something Went Wrong")
		return
	}

	// Respond with JSON array of books
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
