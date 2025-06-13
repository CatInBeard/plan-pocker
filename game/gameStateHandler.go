package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"
)

type RevealRequest struct {
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
}

type StartRequest struct {
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
}

func handleReveal(w http.ResponseWriter, req RevealRequest) {

	gameStateRepository := NewGameStateRepository()
	playerRepository := NewPlayerRepository()

	players, err := playerRepository.GetCachedPlayers(req.GameID)

	if err != nil {
		logger.Log(logger.ERROR, "[GSH-001] Failed to complete vote "+req.GameID, fmt.Sprintf("req: %+v", req))
		http.Error(w, jsonError("Update failed"), http.StatusBadRequest)
		return
	}

	vote := CalculateAveragePositiveVotes(players)

	gameState := GameState{
		ShortLink:  req.GameID,
		VoteStatus: vote,
	}

	gameStateRepository.SetGameState(gameState)

	CalculateGameStateByGameId(req.GameID)

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}

func handleStart(w http.ResponseWriter, req StartRequest) {

	gameStateRepository := NewGameStateRepository()

	gameState := GameState{
		ShortLink:  req.GameID,
		VoteStatus: 0,
	}

	gameStateRepository.SetGameState(gameState)

	CalculateGameStateByGameId(req.GameID)

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}

func CalculateAveragePositiveVotes(players []Player) float64 {
	var totalVotes int
	var positiveVoteCount int

	for _, player := range players {
		if player.Vote > 0 {
			totalVotes += player.Vote
			positiveVoteCount++
		}
	}

	if positiveVoteCount == 0 {
		return 0
	}

	average := float64(totalVotes) / float64(positiveVoteCount)
	return average
}
