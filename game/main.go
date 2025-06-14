package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"shared/httputils"
	"shared/logger"
)

type ActionRequest struct {
	Action string `json:"action"`
}

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		logger.LogFatal(3)
		if r := recover(); r != nil {
			http.Error(w, jsonError("Failed to complete request"), http.StatusInternalServerError)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	requestInfo := httputils.CollectRequestInfoString(r)
	body, _ := io.ReadAll(r.Body)
	logger.Log(logger.DEBUG, "[GHS-003] Received message", fmt.Sprintf("Request info:%s\nBody:%s", requestInfo, body))

	if r.Method != http.MethodPost {
		http.Error(w, jsonError("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	var actionReq ActionRequest
	if err := json.Unmarshal(body, &actionReq); err != nil {
		http.Error(w, jsonError("Invalid JSON"), http.StatusBadRequest)
		return
	}

	switch actionReq.Action {
	case "connect":
		var req ConnectionRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[GHS-004] Error marshal json for ConnectionRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for connection"), http.StatusBadRequest)
			return
		}
		handleConnection(w, req)
	case "updateSettings":
		var req UpdateSettingsRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[GHS-005] Error marshal json for UpdateSettingsRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for updateSettings"), http.StatusBadRequest)
			return
		}
		handleUpdateSettings(w, req)
	case "reveal":
		var req RevealRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[GHS-006] Error marshal json for RevealRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for reveal"), http.StatusBadRequest)
			return
		}
		handleReveal(w, req)
	case "start":
		var req StartRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[GHS-007] Error marshal json for StartRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for start"), http.StatusBadRequest)
			return
		}
		handleStart(w, req)
	case "createGame":
		var req CreateGameRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[GHS-007] Error marshal json for StartRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for start"), http.StatusBadRequest)
			return
		}
		HandleCreateGame(w, req)
	default:
		logger.Log(logger.WARNING, "[GHS-002] Unknown action "+actionReq.Action, fmt.Sprintf("Action: %s, Request info:%s\nBody:%s", actionReq.Action, requestInfo, body))
		http.Error(w, jsonError("Unknown action"), http.StatusBadRequest)
	}
}

func jsonError(message string) string {
	errorResponse := ErrorResponse{Error: message}
	response, _ := json.Marshal(errorResponse)
	return string(response)
}

func main() {
	defer logger.LogFatal()

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8082"
	}
	logger.Log(logger.INFO, "[GHS-001] Http server started", "Http server started on Port "+port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Log(logger.ERROR, "[GHS-002] Failed to serve on port "+port, err.Error())
	}
}
