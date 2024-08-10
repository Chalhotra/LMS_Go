// cmd/pkg/routes/routes.go

package routes

import (
	"bookstore/cmd/pkg/controllers"
	"bookstore/cmd/pkg/middleware"
	"bookstore/cmd/types"
	"bookstore/cmd/views"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Routes for rendering HTML templates
	router.HandleFunc("/register", views.RegisterPage).Methods("GET")
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie := middleware.ExtractTokenFromCookie(r)
		if cookie != "" {
			http.Redirect(w, r, "/api/", http.StatusSeeOther)
		} else {
			views.LoginPage(w, r)
		}
	}).Methods("GET")

	// Routes for handling API requests
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/error", views.RenderErrorPage).Methods("GET")

	// Routes with authentication middleware
	authenticatedRouter := router.PathPrefix("/api").Subrouter()
	authenticatedRouter.Use(middleware.Authenticate)

	authenticatedRouter.HandleFunc("/", redirectToHomeOrLogin).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/welcome", views.AdminHome).Methods("GET")
	authenticatedRouter.HandleFunc("/user/welcome", views.UserHome).Methods("GET")
	authenticatedRouter.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/books/{id}", controllers.GetBook).Methods("GET")
	authenticatedRouter.HandleFunc("/books/{id}", controllers.DeleteBook).Methods("DELETE")
	authenticatedRouter.HandleFunc("/books/{id}", controllers.UpdateBook).Methods("PUT")
	authenticatedRouter.HandleFunc("/search", controllers.SearchBooks).Methods("GET")

	authenticatedRouter.HandleFunc("/logout", controllers.LogoutUser).Methods("POST")
	authenticatedRouter.HandleFunc("/user/request", controllers.RequestAdminStatus).Methods("POST")
	authenticatedRouter.HandleFunc("/user/books/view", views.AvailableBooksPage).Methods("GET")
	authenticatedRouter.HandleFunc("/user/books/available", controllers.GetBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/user/books/checkout/{id}", controllers.CheckoutRequest).Methods("POST")
	authenticatedRouter.HandleFunc("/user/checkouts/history", views.BorrowingHistoryPage).Methods("GET")
	authenticatedRouter.HandleFunc("/user/checkouts/fetch", controllers.FetchCheckouts).Methods("GET")
	authenticatedRouter.HandleFunc("/user/checkins/{checkoutID}", controllers.CheckinBook).Methods("POST")
	authenticatedRouter.HandleFunc("/admin/books", controllers.GetBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/book/{id}", controllers.GetBook).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/books/add", controllers.AdminAddBook).Methods("POST")
	authenticatedRouter.HandleFunc("/admin/books/manage", views.AdminManageBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/books/delete/{id}", controllers.DeleteBook).Methods("DELETE")
	authenticatedRouter.HandleFunc("/admin/books/update/{id}", views.UpdateBookPage).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/books/update/{id}", controllers.UpdateBook).Methods("PUT")
	authenticatedRouter.HandleFunc("/admin/checkouts", views.CheckoutRequestPage).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/checkouts/requests", controllers.GetAllCheckouts).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/checkouts/approve/{id}", controllers.ApproveCheckout).Methods("PUT")
	authenticatedRouter.HandleFunc("/admin/checkouts/deny/{id}", controllers.DenyCheckout).Methods("PUT")
	authenticatedRouter.HandleFunc("/admin/requests", views.AdminRequestsPage).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/requests/view", controllers.GetAdminRequests).Methods("GET")
	authenticatedRouter.HandleFunc("/admin/requests/approve/{id}", controllers.ApproveAdminRequest).Methods("PUT")
	authenticatedRouter.HandleFunc("/admin/requests/deny/{id}", controllers.DenyAdminRequest).Methods("PUT")
	authenticatedRouter.HandleFunc("/admins/view", views.AdminListPage).Methods("GET")
	authenticatedRouter.HandleFunc("/admins", controllers.GetAllAdmins).Methods("GET")
	authenticatedRouter.HandleFunc("/admins/{id}", controllers.RemoveFromAdmin).Methods("DELETE")

	// Add this line for search
	return router

}
func redirectToHomeOrLogin(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(types.User)
	if !ok {
		// If user is not authenticated, redirect to the login page
		fmt.Printf("User not found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.IsAdmin == "1" || user.IsAdmin == "2" {
		// If user is an admin, redirect to the admin welcome page
		http.Redirect(w, r, "/api/admin/welcome", http.StatusSeeOther)
	} else {
		// If user is not an admin, redirect to the user welcome page
		http.Redirect(w, r, "/api/user/welcome", http.StatusSeeOther)
	}
}
