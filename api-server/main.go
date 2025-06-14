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
	logger.Log(logger.DEBUG, "[AS-001] Received message", fmt.Sprintf("Request info:%s\nBody:%s", requestInfo, body))

	if r.Method != http.MethodPost {
		logger.Log(logger.WARNING, "[AS-002] invalid method", fmt.Sprintf("Request info:%s\nBody:%s", requestInfo, body))
		http.Error(w, jsonError("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	var actionReq ActionRequest
	if err := json.Unmarshal(body, &actionReq); err != nil {
		logger.Log(logger.WARNING, "[AS-003] action not found in request", fmt.Sprintf("Request info:%s\nBody:%s", requestInfo, body))
		http.Error(w, jsonError("Invalid JSON"), http.StatusBadRequest)
		return
	}

	switch actionReq.Action {
	case "createGame":
		var req CreateGameRequest
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log(logger.ERROR, "[AS-004] Error marshal json for ConnectionRequest", fmt.Sprintf("Error: %s, Request info:%s\nBody:%s", err, requestInfo, body))
			http.Error(w, jsonError("Invalid JSON for create game"), http.StatusBadRequest)
			return
		}
		handleCreateGameRequest(w, req)
	default:
		logger.Log(logger.WARNING, "[AS-005] Unknown action "+actionReq.Action, fmt.Sprintf("Action: %s, Request info:%s\nBody:%s", actionReq.Action, requestInfo, body))
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
		port = "8083"
	}
	logger.Log(logger.INFO, "[SA-006] Http server started", "Http server started on Port "+port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Log(logger.ERROR, "[AS-007] Failed to serve on port "+port, err.Error())
	}
}
