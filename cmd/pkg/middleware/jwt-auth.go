package middleware

import (
	"bookstore/cmd/pkg/models"
	"bookstore/cmd/types"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Load JWT secret key from environment variable
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Middleware function to authenticate JWT tokens from cookie
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromCookie(r)
		if tokenString == "" {
			// Redirect to login if token is missing
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims := &types.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			// Redirect to login if token is invalid
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !token.Valid {
			// Redirect to login if token is not valid
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := models.GetUserByName(claims.Username)
		if err != nil {
			// Redirect to error page if user not found
			http.Redirect(w, r, "/error?type=404 Not Found&message=User not found", http.StatusSeeOther)
			return
		}

		// Set the user object in the request context
		context.Set(r, "user", user)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Extract JWT token from cookie
func extractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}
