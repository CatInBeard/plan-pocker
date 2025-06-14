package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"shared/logger"
	"sync"
)

var (
	baseURL string
	once    sync.Once
)

type CreateGameApiRequest struct {
	Action string `json:"action"`
	UID    string `json:"uid,omitempty"`
}

type CreateGameApiResponse struct {
	GameID string `json:"gameId"`
}

func getBaseUrl() string {
	once.Do(func() {
		baseURL = os.Getenv("GAME_SERVICE_BASE_URL")
		if baseURL == "" {
			logger.Log(logger.WARNING, "[ASGAC-001] base url not defined", "GAME_SERVICE_BASE_URL not set, use http://game as default")
			baseURL = "http://game"
		}
	})
	return baseURL
}

func SendCreateGameRequest(cg CreateGameApiRequest) (*CreateGameApiResponse, error) {
	cg.Action = "createGame"

	url := getBaseUrl()
	jsonData, err := json.Marshal(cg)

	if err != nil {
		logger.Log(logger.ERROR, "[ASGAC-002] failed to marshal JSON CreateGameRequest", fmt.Sprintf("Error: %s", err.Error()))
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log(logger.ERROR, "[ASGAC-003] failed to send CreateGameRequest", fmt.Sprintf("Error: %s\nRequest body: %s", err.Error(), jsonData))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log(logger.ERROR, "[ASGAC-004] received non-200 response for CreateGameRequest", fmt.Sprintf("Status: %s\nRequest:%s\nResponse:%s", body, resp.Status, jsonData))
		return nil, errors.New("Received " + resp.Status + "status")
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Log(logger.DEBUG, "[ASGAC-005] send CreateGameRequest", fmt.Sprintf("Request: %s, Response: %s", jsonData, body))

	var createGameApiResponse CreateGameApiResponse
	err = json.Unmarshal(body, &createGameApiResponse)
	if err != nil {
		logger.Log(logger.ERROR, "[ASGAC-005] failed to unmarhsal CreateGameApiResponse", fmt.Sprintf("Error: %s\nResponse body: %s", err.Error(), body))
		return nil, err
	}

	return &createGameApiResponse, nil
}
