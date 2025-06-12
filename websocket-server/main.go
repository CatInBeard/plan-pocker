package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"shared/httputils"
	"shared/logger"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientManager struct {
	clients map[string]map[*websocket.Conn]struct{} // gameId -> clients
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (manager *ClientManager) handleConnections(w http.ResponseWriter, r *http.Request) {
	requestInfo := httputils.CollectRequestInfoString(r)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Log(logger.ERROR, "[WS-003] WebSocket upgrade failed", fmt.Sprintf("Error: %s\nRequest details: %s", err.Error(), requestInfo))
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	var gameId string

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			logger.Log(logger.WARNING, "[WS-004] Failed to read websocket message", fmt.Sprintf("Error: %s\nRequest details: %s", err.Error(), requestInfo))
			break
		}
		logger.Log(logger.DEBUG, "[WS-005] Received new WS message", fmt.Sprintf("Message: %s\nHeaders:%s\nRequest details: ", string(message), r.Header, requestInfo))

		var msg map[string]string
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Log(logger.WARNING, "[WS-006] Error unmarshaling message", fmt.Sprintf("Error: %s\nRequest details: %s", err.Error(), requestInfo))
			continue
		}

		if msg["action"] == "reg" {
			gameId = msg["gameId"]
			if manager.clients[gameId] == nil {
				manager.clients[gameId] = make(map[*websocket.Conn]struct{})
			}
			manager.clients[gameId][ws] = struct{}{}
			logger.Log(logger.INFO, fmt.Sprintf("[WS-007] Client registered with gameId: %s", gameId), fmt.Sprintf("Game id: %s\nRequest details: %s", gameId, requestInfo))
		}
	}

	if gameId != "" {
		delete(manager.clients[gameId], ws)
		if len(manager.clients[gameId]) == 0 {
			delete(manager.clients, gameId)
		}
	}
}

func (manager *ClientManager) notifyClients(gameId, message string) {
	if clients, ok := manager.clients[gameId]; ok {
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				logger.Log(logger.WARNING, "[WS-008] Error sending message to client", fmt.Sprintf("Error: %s", err.Error()))
			}
		}
	}
}

func (manager *ClientManager) sendGameNotify(gameId, message string) {
	manager.notifyClients(gameId, message)
}

func main() {
	defer logger.LogFatal()

	manager := NewClientManager()

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", manager.handleConnections)
	logger.Log(logger.INFO, "[WS-001] Http server started", "Http server started on Port "+port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Log(logger.ERROR, "[WS-002] Failed to serve on port "+port, err.Error())
	}
}
