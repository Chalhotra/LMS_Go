package main

import (
	"bookstore/cmd/pkg/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Started the API server")
	r := routes.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
