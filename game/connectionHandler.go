package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"
)

type ConnectionRequest struct {
	Name   string `json:"name"`
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
	Vote   int    `json:"vote,omitempty"`
}

func handleConnection(w http.ResponseWriter, req ConnectionRequest) {

	gameRepository := NewGameRepository()

	_, err := gameRepository.SelectByShortLink(req.GameID)
	if err == sql.ErrNoRows {
		logger.Log(logger.INFO, "[CHI-001] Game with ID "+req.GameID+" not found, creating new gaem", fmt.Sprintf("req: %v", req))
		gameRepository.CreateDefaultGame(req.GameID)
	}

	playerRepository := NewPlayerRepository()
	playerRepository.SetPlayer(Player{
		Name:   req.Name,
		UID:    req.UID,
		Vote:   req.Vote,
		GameId: req.GameID,
	})

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}
