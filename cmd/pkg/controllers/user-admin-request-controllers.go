package controllers

import (
	"bookstore/cmd/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestAdminStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := GetUserFromContext(r)
	if !ok {
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=Unauthorized")
		return
	}
	if user.IsAdmin != "0" {
		json.NewEncoder(w).Encode(map[string]string{"message": "You are already an admin"})
		return
	}
	fmt.Print(user.ID)
	err := models.RequestAdminStatus(user.ID)
	if err != nil {
		if err.Error()[0:41] != "admin request already pending for user ID" {
			w.WriteHeader(http.StatusSeeOther)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
			return
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "U already have a request pending!",
			})
			return

		}

	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Request sent successfully!",
	})

}
