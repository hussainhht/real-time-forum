package main

import (
	"fmt"
	"log"
	"net/http"
	"realtime/backend/db"
	"realtime/backend/config"
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

	fmt.Println("Find the best real time forum on http://localhost" + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Println(err)
		
	}

}
