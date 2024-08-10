package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
		fmt.Print(err.Error())
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
