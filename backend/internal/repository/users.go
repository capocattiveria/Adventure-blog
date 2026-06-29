package repository

import (
	"context"

	"adventure-blog/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles all SQL queries for the users table.
type UserRepository struct {
	pool *pgxpool.Pool // shared connection pool — never copied
}

// NewUserRepository creates a UserRepository with the given connection pool.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user and returns the created row.
// $1, $2 are prepared statement parameters — they prevent SQL injection.
// RETURNING avoids a second SELECT to fetch the generated id and created_at.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.CreatedAt) // & passes pointers so Scan writes directly into the struct fields
	return &u, err
}

// GetByEmail fetches a user by email, including the password hash for login verification.
// Returns an error if no user is found.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return &u, err
}
