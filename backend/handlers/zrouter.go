package handlers

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)

	//! checking method confirming

	// auth
	mux.HandleFunc("POST /register", RegisterHandler) //*dun
	mux.HandleFunc("POST /login", LoginHandler) //*dun
	mux.HandleFunc("POST /logout", LogoutHandler) //* dun

	// user
	mux.HandleFunc("GET /api/session", CurrentUserHandler) //*dun
	mux.HandleFunc("GET /api/chat-users", GetChatUsersHandler) //*dun

	// posts
	mux.HandleFunc("POST /posts", CreatePostHandler) //*dun
	mux.HandleFunc("GET /posts", FeedPostHandler) //*dun
	mux.HandleFunc("GET /posts/{id}", GetPostHandler) //*dun

	// comments
	mux.HandleFunc("POST /posts/{id}/comments", CreateCommentsHandler) //todo on js
	mux.HandleFunc("GET /posts/{id}/comments", FeedCommentsHandler) //todo on js

	// messages
	mux.HandleFunc("GET /api/messages/{userID}", FeedMessagesHandler) //*dun
	mux.HandleFunc("POST /api/messages", MessageHandler) //todo on js

	// websocket
	// mux.HandleFunc("GET /ws", WebSocketHandler) // todo

	return mux
}
