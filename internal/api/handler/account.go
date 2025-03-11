package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"

	"github.com/mark-c-hall/pds-go/internal/api/httputil"
	"github.com/mark-c-hall/pds-go/internal/service"
)

type AccountHandler struct {
	service service.AccountService
	logger  *log.Logger
}

func NewAccountHandler(service service.AccountService, logger *log.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		logger:  logger,
	}
}

type CreateAccountRequest struct {
	Email            string        `json:"email"`
	Handle           syntax.Handle `json:"handle,omitempty"`
	DID              syntax.DID    `json:"did,omitempty"`
	InviteCode       string        `json:"inviteCode,omitempty"`
	VerficationPhone string        `json:"verficationPhone,omitempty"`
	Password         string        `json:"password"`
	RecoveryKey      string        `json:"recoveryKey,omitempty"`
	PLCOperation     struct{}      `json:"plcOp"` // Todo: Figure this out
}

type CreateAccountResponse struct {
	AccessJWT  string               `json:"accessJwt"`
	RefreshJWT string               `json:"refreshJwt"`
	Handle     syntax.Handle        `json:"handle"`
	DID        syntax.DID           `json:"did"`
	DIDDoc     identity.DIDDocument `json:"didDoc"`
}

func (h *AccountHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // TODO: Handle syntax errors
		h.logger.Printf("Failed to decode request: %v", err)
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

	account, err := h.service.CreateAccount(r.Context(), req.Handle, req.Email, req.Password)
	if err != nil {
		h.logger.Printf("Failed to create account: %v", err)
		httputil.RespondWithError(w, "InternalError", "Failed to create account", http.StatusInternalServerError)
		return
	}

	resp := CreateAccountResponse{
		DID:    account.DID,
		Handle: account.Handle,
	}

	h.logger.Printf("Created account for %s with DID %s", account.Handle, account.DID)
	httputil.RespondWithJSON(w, &resp, http.StatusOK)
}
