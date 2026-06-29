package model

import "time"

type Post struct {
	ID           string     `json:"id"`
	AuthorID     string     `json:"author_id"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	Description  *string    `json:"description"`
	ThumbnailURL *string    `json:"thumbnail_url"`
	PublishedAt  *time.Time `json:"published_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
