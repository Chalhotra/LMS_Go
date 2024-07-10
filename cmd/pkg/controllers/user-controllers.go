package controllers

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/pkg/utils"
	"bookstore/cmd/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("soumil05")

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func generateJWT(username string, isAdmin string) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &types.Claims{
		Username: username,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
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

	token, expirationTime, err := generateJWT(tmpuser.Username, tmpuser.IsAdmin)
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

	err = models.RegisterUser(user)
	if err != nil {
		if err.Error() == "username already exists" {
			w.Header().Set("Location", "/error?type=400 Bad Request&message=Username already exists")
		} else {
			w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Internal server error")
		}
		return
	}

	w.Header().Set("Location", "/login")
	json.NewEncoder(w).Encode(map[string]string{"username": user.Username})
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

func CheckoutRequest(w http.ResponseWriter, r *http.Request) {
	// Parse book ID from URL params
	params := mux.Vars(r)
	bookID := params["id"][3:]
	book, err := models.GetBookById(bookID)
	if err != nil {
		w.Header().Set("Location", "/error?type=404 Not Found&message=BookNotFound")

	}

	// Assuming user information is available in context or session
	user, ok := getUserFromContext(r)
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

func getUserFromContext(r *http.Request) (types.User, bool) {
	// Your implementation to retrieve user information from context or session
	// Example implementation assuming user is stored in context
	user, ok := context.Get(r, "user").(types.User)
	return user, ok
}

func FetchCheckouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, ok := getUserFromContext(r)
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

func RequestAdminStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := getUserFromContext(r)
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
