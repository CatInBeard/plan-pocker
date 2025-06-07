package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			return
		}

		fmt.Printf("Received: %s\n", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Failed to write message:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)

	fmt.Println("WebSocket server is running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
