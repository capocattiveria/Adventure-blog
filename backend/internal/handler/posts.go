package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"adventure-blog/internal/model"
)

type PostStore interface {
	ListPublished(ctx context.Context) ([]model.Post, error)
}

type PostHandler struct {
	posts PostStore
}

func NewPostHandler(posts PostStore) *PostHandler {
	return &PostHandler{posts: posts}
}

// List handles GET /posts.
// Public endpoint — returns only published posts.
// Future: when admin auth is in place, admins will also see drafts.
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.posts.ListPublished(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// return an empty array instead of null when there are no posts
	if posts == nil {
		posts = []model.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
