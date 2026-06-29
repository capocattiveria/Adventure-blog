package server

import (
	"net/http"
	"os"

	"adventure-blog/internal/handler"
	custommiddleware "adventure-blog/internal/middleware"
	"adventure-blog/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New builds the HTTP router and wires all dependencies together.
// It receives the DB connection pool and injects it into each handler.
func New(pool *pgxpool.Pool) http.Handler {
	// read allowed origin from env so it can be changed per environment
	// without recompiling (e.g. production frontend URL)
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	r := chi.NewRouter()

	// CORS must be the first middleware so preflight OPTIONS requests
	// are handled before any auth or business logic runs
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// dependency wiring: repository → handler
	userRepo := repository.NewUserRepository(pool)
	authHandler := handler.NewAuthHandler(userRepo)

	// public routes — no authentication required
	r.Get("/health", handler.Health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// protected routes — JWT middleware applied to the whole group
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.Auth)
		// future routes go here
	})

	return r
}
