package handlers

import "net/http"

func Router()*http.ServeMux {
	mux := http.NewServeMux()

	//! checking method confirming

	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)

	// auth
	mux.HandleFunc("POST /register", RegisterHandler)
	mux.HandleFunc("POST /login", LoginHandler)
	mux.HandleFunc("POST /logout", LogoutHandler)

	//user
	mux.HandleFunc("GET /api/session", CurrentUserHandler)

	// posts
	mux.HandleFunc("POST /posts", CreatePostHandler)
	// mux.HandleFunc("GET /posts", FeedHandler)          // todo
	// mux.HandleFunc("GET /posts/{id}", GetPostHandler)  // todo

	// comments
	mux.HandleFunc("POST /posts/{id}/comments", CreateCommentsHandler)
	// mux.HandleFunc("GET /posts/{id}/comments", GetCommentsHandler) // todo

	// messages
	// mux.HandleFunc("GET /messages/{userID}", GetMessagesHandler) // todo
	// mux.HandleFunc("POST /messages", SendMessageHandler)         // todo

	// websocket
	// mux.HandleFunc("GET /ws", WebSocketHandler) // todo

	//todo: WebSocket route
	return mux

}
