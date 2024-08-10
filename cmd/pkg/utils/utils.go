package utils

import (
	"bookstore/cmd/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func GenerateJWT(username string, isAdmin string) (string, time.Time, error) {
	var jwtKey = []byte("soumil05")
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

func GetCurrentUserInfo(w http.ResponseWriter, r *http.Request) types.User {
	user := context.Get(r, "user")

	// Check if user information is present
	if user == nil {
		http.Error(w, "Unauthorized: User information not found", http.StatusUnauthorized)
		return types.User{}
	}
	castedUser, ok := user.(types.User) // Adjust the type based on your user struct
	if !ok {
		http.Error(w, "Internal server error: Unexpected user type in context", http.StatusInternalServerError)
		return types.User{}
	}

	return castedUser

}

func CheckAdmin(w http.ResponseWriter, r *http.Request) bool {
	var user types.User = GetCurrentUserInfo(w, r)

	return user.IsAdmin == "1" || user.IsAdmin == "2"
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")

	}
	return json.NewDecoder(r.Body).Decode(payload)

}

func CheckSuperAdmin(w http.ResponseWriter, r *http.Request) bool {
	var user types.User = GetCurrentUserInfo(w, r)

	return user.IsAdmin == "2"
}
