package handlers

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)

	//! checking method confirming

	// auth
	mux.HandleFunc("POST /register", RegisterHandler)
	mux.HandleFunc("POST /login", LoginHandler)
	mux.HandleFunc("POST /logout", LogoutHandler)

	// user
	mux.HandleFunc("GET /api/session", CurrentUserHandler)
	mux.HandleFunc("GET /api/chat-users", GetChatUsersHandler)

	// posts
	mux.HandleFunc("POST /posts", CreatePostHandler)
	mux.HandleFunc("GET /posts", FeedPostHandler)
	mux.HandleFunc("GET /posts/{id}", GetPostHandler)

	// comments
	mux.HandleFunc("POST /posts/{id}/comments", CreateCommentsHandler)
	mux.HandleFunc("GET /posts/{id}/comments", FeedCommentsHandler)

	// messages
	mux.HandleFunc("GET /messages/{userID}", FeedMessagesHandler)
	mux.HandleFunc("POST /messages", MessageHandler)

	// websocket
	// mux.HandleFunc("GET /ws", WebSocketHandler) // todo

	return mux
}
