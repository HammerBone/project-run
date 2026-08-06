package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	logger      *slog.Logger
	userService *UserService
}

func NewUserHandler(logger *slog.Logger, userService *UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accToken, err := h.userService.CreateUser(r.Context(), &u)
	if err != nil {
		h.logger.Error("failed creating user", slog.Any("error", err))
		http.Error(w, "failed creating user", http.StatusInternalServerError)
		return
	}

	res := UserSignUpRes{
		Name:     u.Name,
		Email:    u.Email,
		AccToken: accToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u UserLoginReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.userService.LoginUser(r.Context(), u.Email, u.Password)
	if err != nil {
		h.logger.Error("error logging in user", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	var u UserLogoutReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.LogoutUser(r.Context(), u.SessionId)
	if err != nil {
		h.logger.Error("error logging out user", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("bad request", slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	newAccToken, err := h.userService.RenewAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("internal server error", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAccToken)
}