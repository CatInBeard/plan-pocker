package main

import (
	"encoding/json"
	"net/http"
	"os"
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
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, jsonError("Invalid request method"), http.StatusMethodNotAllowed)
		return
	}

	var actionReq ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&actionReq); err != nil {
		http.Error(w, jsonError("Invalid JSON"), http.StatusBadRequest)
		return
	}

	switch actionReq.Action {
	case "connection":
		var req ConnectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, jsonError("Invalid JSON for connection"), http.StatusBadRequest)
			return
		}
		handleConnection(w, req)
	case "updateSettings":
		var req UpdateSettingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, jsonError("Invalid JSON for updateSettings"), http.StatusBadRequest)
			return
		}
		handleUpdateSettings(w, req)
	case "reveal":
		var req RevealRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, jsonError("Invalid JSON for reveal"), http.StatusBadRequest)
			return
		}
		handleReveal(w, req)
	case "start":
		var req StartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, jsonError("Invalid JSON for start"), http.StatusBadRequest)
			return
		}
		handleStart(w, req)
	default:
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
