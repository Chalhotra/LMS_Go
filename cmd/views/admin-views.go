package views

import (
	"bookstore/cmd/types"
	"net/http"

	"github.com/gorilla/context"
)

func AdminHome(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "admin-home.html", user)

}

func AdminManageBooks(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "admin-manage-books.html", user)

}

func UpdateBookPage(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "admin-update-book.html", user)

}

func CheckoutRequestPage(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "admin-checkout-requests.html", user)

}

func AdminRequestsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "admin-requests.html", user)

}
