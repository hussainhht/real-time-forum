package handlers

import (
	"net/http"
	"realtime/backend/ws"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	filePath := "./frontend" + r.URL.Path

	// 	_, err := os.Stat(filePath)

	// 	if err == nil {
	// 		http.ServeFile(w, r, filePath)
	// 		return
	// 	}

	// 	http.ServeFile(w, r, "./frontend/index.html")
	// })

	//! checking method confirming

	// auth
	mux.HandleFunc("POST /register", RegisterHandler) //*dun
	mux.HandleFunc("POST /login", LoginHandler)       //*dun
	mux.HandleFunc("POST /logout", LogoutHandler)     //* dun

	// user
	mux.HandleFunc("GET /api/session", CurrentUserHandler)     //*dun
	mux.HandleFunc("GET /api/chat-users", GetChatUsersHandler) //*dun

	// posts
	mux.HandleFunc("POST /posts", CreatePostHandler)  //*dun
	mux.HandleFunc("GET /posts", FeedPostHandler)     //*dun
	mux.HandleFunc("GET /posts/{id}", GetPostHandler) //*dun

	// comments
	mux.HandleFunc("POST /posts/{id}/comments", CreateCommentsHandler) //todo on js
	mux.HandleFunc("GET /posts/{id}/comments", FeedCommentsHandler)    //todo on js

	// messages
	mux.HandleFunc("GET /api/messages/{userID}", FeedMessagesHandler) //*dun
	mux.HandleFunc("POST /api/messages", MessageHandler)              //todo on js

	// websocket
	mux.HandleFunc("GET /ws", ws.WebSocketHandler)

	// mux.Handle("/", http.FileServer(http.Dir("./frontend")))
	mux.HandleFunc("/", serveFrontend)

	return mux
}

var frontendDir = http.Dir("./frontend")

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	f, err := frontendDir.Open(r.URL.Path)
	if err != nil {
		// doesn't exist / traversal attempt rejected -> app shell
		http.ServeFile(w, r, "./frontend/index.html")
		return
	}
	f.Close()

	http.FileServer(frontendDir).ServeHTTP(w, r)
}
