package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"
)

type UpdateSettingsRequest struct {
	UID         string `json:"uid,omitempty"`
	GameID      string `json:"gameId"`
	Deck        []int  `json:"deck"`
	AllowCustom bool   `json:"allowCustom"`
}

func handleUpdateSettings(w http.ResponseWriter, req UpdateSettingsRequest) {

	GameRepository := NewGameRepository()

	err := GameRepository.CreateOrUpdate(Game{
		Shortlink: req.GameID,
		Settings: Settings{
			Deck:            req.Deck,
			AllowCustomDeck: req.AllowCustom,
		},
	})
	if err != nil {
		logger.Log(logger.ERROR, "[USH-001] Failed to update game settings with short link "+req.GameID, fmt.Sprintf("req: %+v", req))
		http.Error(w, jsonError("Update failed"), http.StatusBadRequest)
		return
	}

	logger.Log(logger.DEBUG, "[USH-002] Successfully update game settings for game "+req.GameID, fmt.Sprintf("req: %+v", req))

	CalculateGameStateByGameId(req.GameID)

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}
