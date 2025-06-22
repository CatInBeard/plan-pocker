package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

type GameIdRequest struct {
	UID    string `json:"uid,omitempty"`
	GameId string `json:"gameId"`
}

func (manager *ClientManager) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	requestInfo := httputils.CollectRequestInfoString(r)
	if err != nil {
		requestInfo := httputils.CollectRequestInfoString(r)
		logger.Log(logger.ERROR, "[WS-003] WebSocket upgrade failed", fmt.Sprintf("Error: %s\nRequest details: %s", err.Error(), requestInfo))
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	var gameId string

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				logger.Log(logger.INFO, "[WS-007] Websocket connection closed", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", closeErr.Error(), requestInfo, message))
				return
			}
			logger.Log(logger.WARNING, "[WS-004] Failed to read websocket message", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
			break
		}

		var gameIdRequest GameIdRequest
		if err := json.Unmarshal(message, &gameIdRequest); err != nil {
			logger.Log(logger.WARNING, "[WS-005] Error unmarshaling message to gameId", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
		} else {
			gameId = gameIdRequest.GameId
			if manager.clients[gameId] == nil {
				manager.clients[gameId] = make(map[*websocket.Conn]struct{})
			}
			manager.clients[gameId][ws] = struct{}{}
			logger.Log(logger.DEBUG, fmt.Sprintf("[WS-006] Client registered with gameId: %s", gameId), fmt.Sprintf("Game id: %s\nRequest details: %s\nMessage: %s", gameId, requestInfo, message))
		}

		handleMessage(message, ws, r)
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

func (manager *ClientManager) SendGameNotify(gameId string, message interface{}) {

	stringMessage, err := json.Marshal(message)

	if err != nil {
		logger.Log(logger.WARNING, "[WS-009] Error sending message to clients, failed unmarhal", fmt.Sprintf("Message: %+v\nError: %s", message, err.Error()))
	}

	manager.notifyClients(gameId, string(stringMessage))
}

func main() {
	defer logger.LogFatal()

	manager := NewClientManager()

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", manager.handleConnections)
	go SubscribeGameUpdates(manager.SendGameNotify)
	logger.Log(logger.INFO, "[WS-001] Http server started", "Http server started on Port "+port, logger.Label{Key: "test", Value: "test"})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Log(logger.ERROR, "[WS-002] Failed to serve on port "+port, err.Error())
		time.Sleep(5 * time.Second)
		panic("Failed to start websocket server")
	}
}
