package handlers

import "net/http"

func router() {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/login", Login)
	// mux.HandleFunc("/logout", logout)

	// mux.HandleFunc("/posts",Posts)
	// mux.HandleFunc("/comments",Comments)

	//todo: masge and user hanlders

	//todo: WebSocket route

	fileServer := http.FileServer(http.Dir("./fromtend"))
	mux.Handle("/",fileServer)

}
