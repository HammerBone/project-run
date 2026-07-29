package post

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/HammerBone/project-run/internal/middleware"
)

type PostHandler struct {
	postService *PostService
}

func NewPostHandler(postService *PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
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

	log.Println("[postHandler][CreatePost] claims: ", claims)

	err = h.postService.CreatePost(r.Context(), &p, claims.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
