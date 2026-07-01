package main

import (
	"fmt"
	"log"
	"net/http"
	"realtime/backend/db"
)

func main() {
	//todo: config information
	port := ":8080"

	mux := http.NewServeMux()

	db, err := db.StartDatabase()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	fmt.Println("Find the best real time forum on http://localhost" + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Println(err)
		
	}

}
