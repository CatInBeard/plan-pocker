package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ActionRequest struct {
	Action string `json:"action"`
}

type UpdateSettingsRequest struct {
	Action      string `json:"action"`
	UID         string `json:"uid,omitempty"`
	Deck        []int  `json:"deck"`
	AllowCustom bool   `json:"allowCustom"`
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

func handleUpdateSettings(w http.ResponseWriter, req UpdateSettingsRequest) {
	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}

func handleReveal(w http.ResponseWriter, req RevealRequest) {
	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}

func handleStart(w http.ResponseWriter, req StartRequest) {
	response := Response{Message: "ok"}
	json.NewEncoder(w).Encode(response)
}

func jsonError(message string) string {
	errorResponse := ErrorResponse{Error: message}
	response, _ := json.Marshal(errorResponse)
	return string(response)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
