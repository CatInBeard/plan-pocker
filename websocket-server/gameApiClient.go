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

type ConnectionRequest struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
	Vote   int    `json:"vote,omitempty"`
}

type RevealRequest struct {
	Action string `json:"action"`
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
}

type StartRequest struct {
	Action string `json:"action"`
	UID    string `json:"uid,omitempty"`
	GameID string `json:"gameId"`
}

type UpdateSettingsRequest struct {
	Action      string `json:"action"`
	UID         string `json:"uid,omitempty"`
	GameID      string `json:"gameId"`
	Deck        []int  `json:"deck"`
	AllowCustom bool   `json:"allowCustom"`
}

func getBaseUrl() string {
	once.Do(func() {
		baseURL = os.Getenv("GAME_SERVICE_BASE_URL")
		if baseURL == "" {
			logger.Log(logger.WARNING, "[GAC-001] base url not defined", "GAME_SERVICE_BASE_URL not set, use http://game as default")
			baseURL = "http://game"
		}
	})
	return baseURL
}

func SendConnectRequest(cr ConnectionRequest) error {
	cr.Action = "connect"

	url := getBaseUrl()
	jsonData, err := json.Marshal(cr)

	if err != nil {
		logger.Log(logger.ERROR, "[GAC-002] failed to marshal JSON ConnectionRequest", fmt.Sprintf("Error: %s", err.Error()))
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-003] failed to send ConnectionRequest", fmt.Sprintf("Error: %s\nRequest body: %s", err.Error(), jsonData))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log(logger.ERROR, "[GAC-004] received non-200 response for ConnectionRequest", fmt.Sprintf("Status: %s\nRequest:%s\nResponse:%s", body, resp.Status, jsonData))
		return errors.New("Received " + resp.Status + "status")
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Log(logger.DEBUG, "[GAC-010] send ConnectionRequest", fmt.Sprintf("Request: %s, Response: %s", jsonData, body), logger.Label{Key: "UserUID", Value: cr.UID})

	return nil
}

func SendRevealRequest(rr RevealRequest) error {
	rr.Action = "reveal"

	url := getBaseUrl()
	jsonData, err := json.Marshal(rr)
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-005] failed to marshal JSON RevealRequest", fmt.Sprintf("Error: %s", err.Error()))
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-006] failed to send RevealRequest", fmt.Sprintf("Error: %s\nRequest body: %s", err.Error(), jsonData))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log(logger.ERROR, "[GAC-007] received non-200 response for RevealRequest", fmt.Sprintf("Status: %s\nRequest:%s\nResponse:%s", body, resp.Status, jsonData))
		return errors.New("Received " + resp.Status + "status")
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Log(logger.DEBUG, "[GAC-011] send RevealRequest", fmt.Sprintf("Request: %s, Response: %s", jsonData, body))

	return nil
}

func SendStartRequest(sr StartRequest) error {
	sr.Action = "start"

	url := getBaseUrl()
	jsonData, err := json.Marshal(sr)
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-008] failed to marshal JSON StartRequest", fmt.Sprintf("Error: %s\nRequest body: %s", err.Error(), jsonData))
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-009] failed to send StartRequest", fmt.Sprintf("Error: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log(logger.ERROR, "[GAC-010] received non-200 response for StartRequest", fmt.Sprintf("Status: %s\nRequest:%s\nResponse:%s", body, resp.Status, jsonData))
		return errors.New("Received " + resp.Status + "status")
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Log(logger.DEBUG, "[GAC-011] send StartRequest", fmt.Sprintf("Request: %s, Response: %s", jsonData, body))

	return nil
}

func SendUpdateSettingsRequest(ur UpdateSettingsRequest) error {
	ur.Action = "updateSettings"

	url := getBaseUrl()
	jsonData, err := json.Marshal(ur)
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-011] failed to marshal JSON UpdateSettingsRequest", fmt.Sprintf("Error: %s\nRequest body: %s", err.Error(), jsonData))
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log(logger.ERROR, "[GAC-012] failed to send UpdateSettingsRequest", fmt.Sprintf("Error: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log(logger.ERROR, "[GAC-013] received non-200 response for UpdateSettingsRequest", fmt.Sprintf("Status: %s\nRequest:%s\nResponse:%s", body, resp.Status, jsonData))
		return errors.New("Received " + resp.Status + "status")
	}

	body, _ := io.ReadAll(resp.Body)
	logger.Log(logger.DEBUG, "[GAC-011] send UpdateSettingsRequest", fmt.Sprintf("Request: %s, Response: %s", jsonData, body))

	return nil
}

func GetConnectionRequestFromWsRequest(wsConnectionRequest WsConnectionRequest) ConnectionRequest {
	return ConnectionRequest{
		Action: "connect",
		Name:   wsConnectionRequest.UserName,
		UID:    wsConnectionRequest.UID,
		GameID: wsConnectionRequest.GameId,
		Vote:   wsConnectionRequest.Vote,
	}
}

func GetRevealRequestFromWsRequest(wsConnectionRequest WsRevealRequest) RevealRequest {
	return RevealRequest{
		Action: "reveal",
		UID:    wsConnectionRequest.UID,
		GameID: wsConnectionRequest.GameId,
	}
}

func GetStartRequestFromWsRequest(wsConnectionRequest WsStartRequest) StartRequest {
	return StartRequest{
		Action: "start",
		UID:    wsConnectionRequest.UID,
		GameID: wsConnectionRequest.GameId,
	}
}

func GetUpdateSettingsRequestFromWsRequest(wsConnectionRequest WsSettingsRequest) UpdateSettingsRequest {
	return UpdateSettingsRequest{
		Action:      "updateSettings",
		UID:         wsConnectionRequest.UID,
		GameID:      wsConnectionRequest.GameId,
		Deck:        wsConnectionRequest.Deck,
		AllowCustom: wsConnectionRequest.AllowCustom,
	}
}
