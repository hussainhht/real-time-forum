package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintln(w, "Hello, World!")
	})

	port := ":8080"

	fmt.Println("Find the best real time forum on http://localhost" + port)
	http.ListenAndServe(port, mux)

}
