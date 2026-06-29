package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"adventure-blog/internal/handler"
	"adventure-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type mockUserStore struct {
	createFn     func(ctx context.Context, email, hash string) (*model.User, error)
	getByEmailFn func(ctx context.Context, email string) (*model.User, error)
}

func (m *mockUserStore) Create(ctx context.Context, email, hash string) (*model.User, error) {
	return m.createFn(ctx, email, hash)
}

func (m *mockUserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.getByEmailFn(ctx, email)
}

func TestRegister_Created(t *testing.T) {
	store := &mockUserStore{
		createFn: func(ctx context.Context, email, hash string) (*model.User, error) {
			return &model.User{ID: "123", Email: email}, nil
		},
	}

	h := handler.NewAuthHandler(store)
	body := `{"email":"test@test.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var user model.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != "test@test.com" {
		t.Errorf("unexpected email: %s", user.Email)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	h := handler.NewAuthHandler(&mockUserStore{})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegister_EmailTaken(t *testing.T) {
	store := &mockUserStore{
		createFn: func(ctx context.Context, email, hash string) (*model.User, error) {
			return nil, errors.New("duplicate key")
		},
	}

	h := handler.NewAuthHandler(store)
	body := `{"email":"taken@test.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)

	store := &mockUserStore{
		getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{ID: "123", Email: email, PasswordHash: string(hash)}, nil
		},
	}

	h := handler.NewAuthHandler(store)
	body := `{"email":"test@test.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)

	store := &mockUserStore{
		getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{ID: "123", Email: email, PasswordHash: string(hash)}, nil
		},
	}

	h := handler.NewAuthHandler(store)
	body := `{"email":"test@test.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	store := &mockUserStore{
		getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	h := handler.NewAuthHandler(store)
	body := `{"email":"nobody@test.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
