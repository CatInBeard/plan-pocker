package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shared/cache"
	"shared/logger"
	"strconv"
	"time"
)

type ConnectionRequest struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
	Vote   int    `json:"vote,omitempty"`
}

type Player struct {
	Name   string
	UID    string
	GameId string
	Vote   int
}

func handleConnection(w http.ResponseWriter, req ConnectionRequest) {

	cacheLiveDurationString := GetSetting(STAY_CONNECTED_PLAYER_SETTING)
	cacheLiveDuration, _ := strconv.Atoi(cacheLiveDurationString)

	GameRepository := NewGameRepository()

	_, err := GameRepository.SelectByShortLink(req.GameID)
	if err == sql.ErrNoRows {
		logger.Log(logger.INFO, "[CHI-001] Game with ID "+req.GameID+" not found, creating new gaem", fmt.Sprintf("req: %v", req))
		GameRepository.CreateDefaultGame(req.GameID)
	}

	c := cache.GetCacheClient()
	c.SetStructValue(
		fmt.Sprintf("game_%s_player", req.GameID),
		Player{
			Name:   req.Name,
			UID:    req.UID,
			GameId: req.GameID,
			Vote:   req.Vote,
		},
		time.Duration(cacheLiveDuration)*time.Second)

	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}
