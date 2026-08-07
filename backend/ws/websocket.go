package ws

import (
	"log"
	"net/http"
	"realtime/backend/global/utilities"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "failed to upgrade to websocket")
		log.Println("Failed to upgrade to websocket:", err)
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		//messageType this for text or binary ,, binary is for file and text is for text
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("Received message: %s", message) //if you want to know waht is the massage

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			utilities.ErrorJSON(w, http.StatusInternalServerError, "failed to write message to websocket")
			log.Println("Error writing message:", err)
			return
		}

	}

}
