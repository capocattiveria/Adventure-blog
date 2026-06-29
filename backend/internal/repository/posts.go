package repository

import (
	"context"

	"adventure-blog/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

// ListPublished returns all published posts ordered from newest to oldest.
// content is intentionally excluded — it is only needed on the single-post page.
func (r *PostRepository) ListPublished(ctx context.Context) ([]model.Post, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, author_id, title, slug, description, thumbnail_url, published_at, created_at, updated_at
		FROM posts
		WHERE published_at IS NOT NULL
		ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.ID, &p.AuthorID, &p.Title, &p.Slug,
			&p.Description, &p.ThumbnailURL, &p.PublishedAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
