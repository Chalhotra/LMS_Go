package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AdminAddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book types.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		// w.WriteHeader(http.StatusSeeOther)
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Invalid request format")
		return
	}

	err = models.AddBook(book)
	if err != nil {

		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return
	}

	// w.WriteHeader(http.StatusSeeOther)
	// w.Header().Set("Location", "/login")
	json.NewEncoder(w).Encode(map[string]string{"message": "Book added successfully"})
}

func GetAllCheckouts(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an admin
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")
		return
	}

	// Fetch all checkout requests from the database
	checkouts, err := models.GetAllCheckouts()
	if err != nil {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")

		return
	}

	// Convert checkouts to JSON and send response
	json.NewEncoder(w).Encode(checkouts)
}

func ApproveCheckout(w http.ResponseWriter, r *http.Request) {

	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")

		return
	}
	params := mux.Vars(r)
	checkoutID := params["id"]

	err := models.ApproveCheckout(checkoutID)

	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Approved Checkout Request!"})
}

func DenyCheckout(w http.ResponseWriter, r *http.Request) {

	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")

		return
	}
	params := mux.Vars(r)
	checkoutID := params["id"]

	err := models.DenyCheckout(checkoutID)

	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Denied Checkout Request!"})
}

func GetAdminRequests(w http.ResponseWriter, r *http.Request) {

	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")
		return
	}

	adminRequests, err := models.GetAdminRequests()

	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return
	}

	json.NewEncoder(w).Encode(adminRequests)

}

func ApproveAdminRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")

		return
	}

	params := mux.Vars(r)
	userID := params["id"]
	err := models.ApproveAdminRequest(userID)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error!",
		})

		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return

	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Approved Admin Request!",
	})

}

func DenyAdminRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !utils.CheckAdmin(w, r) {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=You are not an admin")

		return
	}

	params := mux.Vars(r)
	userID := params["id"]
	err := models.DenyAdminRequest(userID)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error",
		})
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		return

	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Denied Admin Request!",
	})

}
func GetAllAdmins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is a super admin
	if !utils.CheckSuperAdmin(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	admins, err := models.GetAllAdmins()
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(admins)
}
func RemoveFromAdmin(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckSuperAdmin(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	params := mux.Vars(r)
	userID := params["id"]
	err := models.RemoveFromAdmin(userID)
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Removed from admin successfully"})

}
