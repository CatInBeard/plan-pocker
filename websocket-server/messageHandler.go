package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"shared/httputils"
	"shared/logger"
)

type ActionRequest struct {
	Action string `json:"action"`
}

type WsConnectionRequest struct {
	UserName string `json:"userName"`
	UID      string `json:"uid,omitempty"`
	GameId   string `json:"gameId"`
	Vote     int    `json:"vote"`
}

type WsRevealRequest struct {
	UID    string `json:"uid,omitempty"`
	GameId string `json:"gameId"`
}

type WsStartRequest struct {
	UID    string `json:"uid,omitempty"`
	GameId string `json:"gameId"`
}

type WsSettingsRequest struct {
	UID         string `json:"uid,omitempty"`
	GameId      string `json:"gameId"`
	AllowCustom bool   `json:"isCustomDeckAllowed"`
	Deck        []int  `json:"deck"`
}

func handleMessage(message []byte, ws *websocket.Conn, r *http.Request) {
	requestInfo := httputils.CollectRequestInfoString(r)
	logger.Log(logger.DEBUG, "[WSH-001] Received new WS message", fmt.Sprintf("Message: %s\nHeaders:%s\nRequest details: %s", string(message), r.Header, requestInfo))

	var actionRequest ActionRequest
	if err := json.Unmarshal(message, &actionRequest); err != nil {
		logger.Log(logger.WARNING, "[WSH-002] Error unmarshaling message to get action", fmt.Sprintf("Error: %s\nRequest details: %s", err.Error(), requestInfo))
		return
	}

	switch actionRequest.Action {
	case "connect":
		var wsConnectionRequest WsConnectionRequest
		if err := json.Unmarshal(message, &wsConnectionRequest); err != nil {
			logger.Log(logger.WARNING, "[WSH-003] Error unmarshaling messaWsConnectionRequest", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
			return
		}
		SendConnectRequest(GetConnectionRequestFromWsRequest(wsConnectionRequest))

		break
	case "reveal":
		var wsRevealRequest WsRevealRequest
		if err := json.Unmarshal(message, &wsRevealRequest); err != nil {
			logger.Log(logger.WARNING, "[WSH-004] Error unmarshaling WsRevealRequest", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
			return
		}
		SendRevealRequest(GetRevealRequestFromWsRequest(wsRevealRequest))
		break
	case "start":
		var wsStartRequest WsStartRequest
		if err := json.Unmarshal(message, &wsStartRequest); err != nil {
			logger.Log(logger.WARNING, "[WSH-005] Error unmarshaling WsStartRequest", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
			return
		}
		SendStartRequest(GetStartRequestFromWsRequest(wsStartRequest))
		break
	case "setSettings":
		var wsSettingsRequest WsSettingsRequest
		if err := json.Unmarshal(message, &wsSettingsRequest); err != nil {
			logger.Log(logger.WARNING, "[WSH-006] Error unmarshaling WsSettingsRequest", fmt.Sprintf("Error: %s\nRequest details: %s\nMessage: %s", err.Error(), requestInfo, message))
			return
		}
		SendUpdateSettingsRequest(GetUpdateSettingsRequestFromWsRequest(wsSettingsRequest))
		break
	default:
		logger.Log(logger.WARNING, "[WSH-007] Invalid action "+actionRequest.Action, fmt.Sprintf("Action: %s\nRequest details: %s\nMessage: %s", actionRequest.Action, requestInfo, message))
	}

}
