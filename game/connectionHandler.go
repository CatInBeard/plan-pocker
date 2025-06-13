package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"

	"github.com/go-redis/redis/v8"
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
		logger.Log(logger.INFO, "[CHI-001] Game with ID "+req.GameID+" not found, creating new gaem", fmt.Sprintf("req: %+v", req))
		gameRepository.CreateDefaultGame(req.GameID)
	}

	gameStateRepository := NewGameStateRepository()

	gameState, errState := gameStateRepository.GetGameState(req.GameID)

	playerRepository := NewPlayerRepository()

	shouldUpdate := false

	if errState != nil {
		shouldUpdate = true
	} else {
		if gameState.VoteStatus == 0 {
			shouldUpdate = true
		} else {
			_, err = playerRepository.GetPlayer(req.GameID, req.UID)
			if err == redis.Nil {
				shouldUpdate = true
			}
		}
	}

	if shouldUpdate {
		playerRepository.SetPlayer(Player{
			Name:   req.Name,
			UID:    req.UID,
			Vote:   req.Vote,
			GameId: req.GameID,
		})
	} else {
		playerRepository.UpdatePlayerLive(Player{
			UID:    req.UID,
			GameId: req.GameID,
		})
	}

	go CalculateGameStateByGameId(req.GameID)

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}
