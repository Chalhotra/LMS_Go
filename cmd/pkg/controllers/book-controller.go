package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an admin
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=User is not an admin")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var body types.Book
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Invalid request format")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	err = models.AddBook(body)
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an admin
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=User is not an admin")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	err := models.DeleteBook(id)
	if err != nil {

		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Book Deleted successfully"})

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an admin
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=User is not an admin")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book types.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Invalid request format")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	err = models.UpdateBook(params["id"], book)
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

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
	bookList, err := models.GetBooks()
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
