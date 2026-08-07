package post

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/HammerBone/project-run/internal/middleware"
)

type PostHandler struct {
	logger      *slog.Logger
	postService *PostService
}

func NewPostHandler(logger *slog.Logger, postService *PostService) *PostHandler {
	return &PostHandler{
		logger:      logger,
		postService: postService,
	}
}

func (h *PostHandler) GetAllPost(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	res, err := h.postService.GetAllPost(r.Context(), claims.Id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var p Post

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	err = h.postService.CreatePost(r.Context(), &p, claims.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	var p Post

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	res, err := h.postService.EditPost(r.Context(), &p, claims.Id)
	if err != nil {
		h.logger.Error("failed to edit post", slog.Any("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
