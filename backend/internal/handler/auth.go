package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"adventure-blog/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserStore defines the DB operations required for authentication.
// Using an interface instead of a concrete struct allows swapping
// the implementation in tests without touching a real database.
type UserStore interface {
	Create(ctx context.Context, email, passwordHash string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	users UserStore
}

// NewAuthHandler creates an AuthHandler with the injected dependency.
func NewAuthHandler(users UserStore) *AuthHandler {
	return &AuthHandler{users: users}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register creates a new user.
// POST /auth/register — body: { email, password }
// Returns 201 with user data, 409 if the email is already taken.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// bcrypt includes an automatic random salt — two users with the same
	// password will have different hashes. DefaultCost balances security and speed.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := h.users.Create(r.Context(), req.Email, string(hash))
	if err != nil {
		// the DB returns an error if the email violates the UNIQUE constraint
		http.Error(w, "email already taken", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login verifies credentials and returns a JWT.
// POST /auth/login — body: { email, password }
// Returns 200 with { token }, 401 if credentials are invalid.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		// return the same error for both email not found and wrong password
		// to avoid revealing which one is incorrect.
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(user)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// generateJWT creates a JWT signed with HS256.
// The token contains the user ID (sub), email, and expires after 24 hours.
// The signing key is read from JWT_SECRET — in production it must be
// a long random string kept secret.
func generateJWT(user *model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
