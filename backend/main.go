package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"realtime/db"
)

func main() {
	mux := http.NewServeMux()

	db ,err := db.StartDatabase()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintln(w, "Hello, World!")
	})

	port := ":8080"

	fmt.Println("Find the best real time forum on http://localhost" + port)
	http.ListenAndServe(port, mux)


}
