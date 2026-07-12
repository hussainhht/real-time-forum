package main

import (
	"fmt"
	"log"
	"net/http"
	"realtime/backend/config"
	"realtime/backend/db"
	"realtime/backend/handlers"
)

func main() {

	cfg := config.Load();

	port := cfg.Port
	DBPath := cfg.DBPath

	mux := http.NewServeMux()

	db, err := db.StartDatabase(DBPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/register", handlers.RegisterHandler) //? write handler in the name or not?

	fmt.Println("Find the best real time forum on http://localhost" + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Println(err)
	}

}
