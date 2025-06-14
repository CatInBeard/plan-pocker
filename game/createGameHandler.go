package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"
	"time"
)

type CreateGameRequest struct {
	Action string `json:"action"`
	UID    string `json:"uid,omitempty"`
}

type CreateGameResponse struct {
	GameID string `json:"gameId"`
}

func HandleCreateGame(w http.ResponseWriter, cgr CreateGameRequest) {
	gameRepository := NewGameRepository()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			logger.Log(logger.ERROR, "[CGH-003] Failed to create a game", fmt.Sprintf("Error: %s, Request info:%+v", "timeout", cgr))
			http.Error(w, jsonError("Failed to create a game"), http.StatusTeapot)
			return
		default:
			shortLink := GenerateShortString(10)
			_, err := gameRepository.SelectByShortLink(shortLink)

			if err == sql.ErrNoRows {
				logger.Log(logger.INFO, "[CGH-001] Successfully create a game", fmt.Sprintf("Request info:%+v", cgr))
				err = gameRepository.CreateDefaultGame(shortLink)
				if err != nil {
					logger.Log(logger.ERROR, "[CGH-004] Error to create a game", fmt.Sprintf("Error: %s, Request info:%+v", err, cgr))
					http.Error(w, jsonError("Failed to create a game"), http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(CreateGameResponse{GameID: shortLink})
				return
			} else if err != nil {
				logger.Log(logger.ERROR, "[CGH-002] Error to create a game", fmt.Sprintf("Error: %s, Request info:%+v", err, cgr))
				http.Error(w, jsonError("Failed to create a game"), http.StatusInternalServerError)
				return
			}
		}
	}

}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
