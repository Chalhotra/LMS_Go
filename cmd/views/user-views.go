// cmd/views/user-views.go

package views

import (
	"bookstore/cmd/types"

	"html/template"
	"net/http"

	"github.com/gorilla/context"
)

// RenderTemplate is a helper function to render HTML templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	user, _ := context.Get(r, "user").(types.User)

	w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
	tmplPath := "cmd/templates/" + tmpl
	var t *template.Template
	var err error
	if user.IsAdmin == "1" {
		t, err = template.ParseFiles(tmplPath, "cmd/templates/nav-admin.html")
	} else {
		t, err = template.ParseFiles(tmplPath, "cmd/templates/nav.html")
	}

	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	t.Execute(w, data)
}

// RegisterPage renders the registration page
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "register.html", nil)
}

// LoginPage renders the login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "login.html", nil)
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Render the template with user data
	RenderTemplate(w, r, "user-home.html", user)

}

func AvailableBooksPage(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	RenderTemplate(w, r, "view-books.html", user)
}

func BorrowingHistoryPage(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=User Not Found in Context")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	RenderTemplate(w, r, "borrowing-history.html", user)
}
