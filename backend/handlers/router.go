package handlers

import "net/http"

func Router()*http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("/register", RegisterHandler) //? write handler in the name or not?
	mux.HandleFunc("/login", LoginHandler) //* : "POST/login" try to do this for handle the path and login 
	// mux.HandleFunc("/logout", logout)

	mux.HandleFunc("/create-post", CreatePostHandler)
	//feed handler

	mux.HandleFunc("/create-comment", CreateCommentsHandler)
	//get comments handler

	//todo: masge and user hanlders

	//todo: WebSocket route
	return mux

}
