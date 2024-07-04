package views

import (
	"html/template"
	"net/http"
)

type ErrorData struct {
	Type    string
	Message string
}

func RenderErrorPage(w http.ResponseWriter, r *http.Request) {
	errorType := r.URL.Query().Get("type")
	if errorType == "" {
		errorType = "Unknown Error"
	}

	errorMessage := r.URL.Query().Get("message")
	if errorMessage == "" {
		errorMessage = "An unknown error occurred."
	}

	tmplPath := "cmd/templates/errors.html"
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Something Went Wrong")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	data := ErrorData{
		Type:    errorType,
		Message: errorMessage,
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {

		w.Header().Set("Location", "/error?type=500 Internal Server Error&message=Something Went Wrong")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
