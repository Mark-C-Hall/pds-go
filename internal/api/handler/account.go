package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"

	"github.com/mark-c-hall/pds-go/internal/api/httputil"
)

type CreateAccountRequest struct {
	Email            string        `json:"email"`
	Handle           syntax.Handle `json:"handle,omitempty"`
	DID              syntax.DID    `json:"did,omitempty"`
	InviteCode       string        `json:"inviteCode,omitempty"`
	VerficationPhone string        `json:"verficationPhone,omitempty"`
	Password         string        `json:"password"`
	RecoveryKey      string        `json:"recoveryKey,omitempty"`
	PLCOperation     struct{}      `json:"plcOp,omitempty"` // Todo: Figure this out
}

type CreateAccountResponse struct {
	AccessJWT  string               `json:"accessJwt"`
	RefreshJWT string               `json:"refreshJwt"`
	Handle     syntax.Handle        `json:"handle"`
	DID        syntax.DID           `json:"did"`
	DIDDoc     identity.DIDDocument `json:"didDoc"`
}

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, "InvalidRequest", "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Handle == "" {
		httputil.RespondWithError(w, "InvalidRequest", "Handle is required", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		httputil.RespondWithError(w, "InvalidRequest", "Email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		httputil.RespondWithError(w, "InvalidRequest", "Password is required", http.StatusBadRequest)
		return
	}

	resp := CreateAccountResponse{Handle: req.Handle}
	log.Printf("Created account for %s\n", req.Handle)
	httputil.RespondWithJSON(w, &resp, http.StatusOK)
}
