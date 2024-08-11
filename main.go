package main

import (
	"caapp-server/src/database"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// init db
	database.Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the server!")
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
