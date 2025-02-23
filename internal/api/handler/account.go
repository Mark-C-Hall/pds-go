package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type CreateAccountRequest struct {
	Email            string   `json:"email"`
	Handle           string   `json:"handle"`
	DID              string   `json:"did"` // Todo: change type later
	InviteCode       string   `json:"inviteCode"`
	VerficationPhone string   `json:"verficationPhone"`
	Password         string   `json:"password"`
	RecoveryKey      string   `json:"recoveryKey"`
	PLCOperation     struct{} `json:"plcOp"` // Todo: Figure this out
}

type CreateAccountResponse struct {
	AccessJWT  string   `json:"accessJwt"`
	RefreshJWT string   `json:"refreshJwt"`
	Handle     string   `json:"handle"`
	DID        string   `json:"did"`    // Todo: Fix type
	DIDDoc     struct{} `json:"didDoc"` // Todo: This
}

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, "InvalidRequest", "Method Not Supported", http.StatusBadRequest)
		return
	}

	if !slices.Contains(r.Header["Content-Type"], "application/json") {
		respondWithError(w, "InvalidRequest", "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "InvalidRequest", "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Handle == "" {
		respondWithError(w, "InvalidRequest", "Handle is required", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		respondWithError(w, "InvalidRequest", "Email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		respondWithError(w, "InvalidRequest", "Password is required", http.StatusBadRequest)
		return
	}

	resp := CreateAccountResponse{Handle: req.Handle}
	log.Printf("Created account for %s\n", req.Handle)
	respondWithJSON(w, &resp, http.StatusOK)
}

func respondWithError(w http.ResponseWriter, err, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err,
		Message: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, res *CreateAccountResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(*res); err != nil {
		log.Printf("Failed to encode JSON response: %v\n", err)
	}
}
