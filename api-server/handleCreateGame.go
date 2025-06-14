package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shared/logger"
)

type CreateGameRequest struct {
	UID string `json:"uid,omitempty"`
}

func handleCreateGameRequest(w http.ResponseWriter, cgr CreateGameRequest) {
	createGameApiResponse, err := SendCreateGameRequest(CreateGameApiRequest{
		UID: cgr.UID,
	})

	if err != nil {
		logger.Log(logger.ERROR, "[ASCH-001] Failed to create game", fmt.Sprintf("Request info:%+v\nError:%s", cgr, err.Error()))
		http.Error(w, jsonError("Invalid JSON"), http.StatusBadRequest)
		return
	}

	logger.Log(logger.INFO, "[ASCH-002] Game created", fmt.Sprintf("Request info:%+v\nResponse:%s", cgr, createGameApiResponse))
	json.NewEncoder(w).Encode(createGameApiResponse)
}
