package handlers

import "net/http"

func Router() *http.ServeMux {
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
	mux.HandleFunc("GET /api/chat-users", GetChatUsersHandler)	//users get users handler

	// posts
	mux.HandleFunc("POST /posts", CreatePostHandler)
	mux.HandleFunc("GET /posts", FeedPostHandler)          // todo
	// mux.HandleFunc("GET /posts/{id}", GetPostHandler)  // todo


	mux.HandleFunc("POST /api/posts/{id}/like", LikeHandler)
	mux.HandleFunc("POST /api/posts/{id}/dislike", DislikeHandler)

	// comments
	mux.HandleFunc("POST /posts/{id}/comments", CreateCommentsHandler)
	mux.HandleFunc("GET /posts/{id}/comments", FeedCommentsHandler) 

	// messages
	mux.HandleFunc("GET /messages/{userID}", GetMessagesHandler) // todo
	mux.HandleFunc("POST /messages", MessageHandler)         // todo

	// websocket
	// mux.HandleFunc("GET /ws", WebSocketHandler) // todo

	//todo: WebSocket route
	return mux

}
