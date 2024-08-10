package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user types.JsonResponse
	w.Header().Set("Content-Type", "application/json")
	var tmpuser types.User
	err := utils.ParseJSON(r, &user)
	if err != nil {
		// w.WriteHeader(http.StatusSeeOther)
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Invalid request format")
		return
	}

	tmpuser, err = models.GetUserByName(user.Username)
	if err != nil {
		// w.WriteHeader(http.StatusSeeOther)
		w.Header().Set("Location", "/error?type=404 Not Found&message=User not found")
		return
	}

	if !checkPassword(tmpuser.Password, user.Password) {
		// w.WriteHeader(http.StatusSeeOther)
		w.Header().Set("Location", "/error?type=401 Unauthorized&message=Credentials do not match")
		return
	}

	token, expirationTime, err := utils.GenerateJWT(tmpuser.Username, tmpuser.IsAdmin)
	if err != nil {
		// w.WriteHeader(http.StatusSeeOther)
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Failed to sign JWT")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: expirationTime,
	})

	// w.WriteHeader(http.StatusSeeOther)
	if tmpuser.IsAdmin == "1" || tmpuser.IsAdmin == "2" {
		w.Header().Set("Location", "/api/admin/welcome")
	} else {
		w.Header().Set("Location", "/api/user/welcome")
	}

}
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Location", "/error?type=400 Bad Request&message=Invalid request format")
		return
	}

	if len(user.Password) >= 8 {
		err = models.RegisterUser(user)
		if err != nil {
			if err.Error() == "username already exists" {
				w.Header().Set("Location", "/error?type=400 Bad Request&message=Username already exists")

			} else {
				w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
			}
			return
		}

		// w.Header().Set("Location", "/login")
		// json.NewEncoder(w).Encode(map[string]string{"username": user.Username})

		token, expirationTime, err := utils.GenerateJWT(user.Username, user.IsAdmin)
		if err != nil {
			// w.WriteHeader(http.StatusSeeOther)
			w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Failed to sign JWT")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   token,
			Expires: expirationTime,
		})

		// w.WriteHeader(http.StatusSeeOther)

		w.Header().Set("Location", "/api/user/welcome")

	} else {
		w.Header().Set("Location", "/error?type=400 Bad Request Error&message=Password must be at least 8 characters")

	}

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Expire the JWT cookie by setting its expiration time to a time in the past
	expiration := time.Now().AddDate(0, 0, -1) // Set expiration to 1 day ago
	cookie := http.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	// Optionally, redirect to a different page after logout
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusSeeOther)
}

// cmd/pkg/controllers/checkout.go
