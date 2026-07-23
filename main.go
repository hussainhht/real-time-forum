package main

import (
	"fmt"
	"log"
	"net/http"
	"realtime/backend/config"
	"realtime/backend/db"
	"realtime/backend/global"
	"realtime/backend/handlers"
)

func main() {

	cfg := config.Load()
	port := cfg.Port
	DBPath := cfg.DBPath
	mux := handlers.Router()

	db, err := db.StartDatabase(DBPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	global.Database = db

	fmt.Println("Find the best real time forum on http://localhost" + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Println(err)
	}

}
