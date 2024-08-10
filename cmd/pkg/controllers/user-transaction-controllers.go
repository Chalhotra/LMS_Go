package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func CheckoutRequest(w http.ResponseWriter, r *http.Request) {
	// Parse book ID from URL params
	params := mux.Vars(r)
	bookID := params["id"][3:]
	book, err := models.GetBookById(bookID)
	if err != nil {
		w.Header().Set("Location", "/error?type=404 Not Found&message=BookNotFound")

	}

	// Assuming user information is available in context or session
	user, ok := GetUserFromContext(r)
	if !ok {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=Unauthorized")
		return
	}

	// Calculate due date (e.g., 14 days from checkout date)
	checkoutDate := time.Now()
	dueDate := checkoutDate.AddDate(0, 0, 14)

	// Create a checkout request
	checkout := types.CheckoutRequest{
		UserID:       user.ID,
		BookID:       book.ID,
		CheckoutDate: checkoutDate,
		DueDate:      dueDate,
	}

	// Add the checkout request to the database
	err = models.CreateCheckoutRequest(checkout)

	if err != nil {
		if err.Error()[0:25] == "user has already borrowed" {
			w.Header().Set("Location", "/error?type=400 Bad Request&message=You have already borrowed this book right now!")
			return

		} else {
			w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Failed to create checkout")
			return
		}
	}

	// Respond with success message or JSON response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkout)
}

func GetUserFromContext(r *http.Request) (types.User, bool) {
	// Your implementation to retrieve user information from context or session
	// Example implementation assuming user is stored in context
	user, ok := context.Get(r, "user").(types.User)
	return user, ok
}

func FetchCheckouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, ok := GetUserFromContext(r)
	if !ok {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=Unauthorized")
		return
	}
	bookList, err := models.GetCheckoutHistory(user.ID)
	// bookList, err := models.GetBooks()
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(bookList)
}
func CheckinBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	checkoutID, err := strconv.Atoi(params["checkoutID"])
	if err != nil {

		log.Printf("Invalid checkout ID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid checkout ID",
		})
		return
	}

	err = models.CheckinBook(checkoutID)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to check in book",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book checked in successfully",
	})
}
