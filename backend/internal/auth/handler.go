package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	service      ServiceInterface
	tokenManager TokenManagerInterface
}

func NewHandler(service ServiceInterface, tkManager TokenManagerInterface) *Handler {
	return &Handler{service: service, tokenManager: tkManager}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := registerResponse{user.Id, user.Username, user.IsActive, user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	resp := LoginResponse{res.user.Id, res.user.Username, res.AccessToken, res.ExpiresAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	parts := strings.Fields(r.Header.Get("Authorization"))

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := h.tokenManager.VerifyAccessToken(parts[1])
	if err != nil {
		http.Error(w, "bad jwt", http.StatusForbidden)
		return
	}

	user, err := h.service.GetUserInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, "unablel to find user", http.StatusBadGateway)
		return
	}

	resp := registerResponse{user.Id, user.Username, user.IsActive, user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
