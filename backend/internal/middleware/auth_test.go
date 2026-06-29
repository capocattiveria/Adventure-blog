package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"adventure-blog/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

// validToken genera un JWT firmato valido da usare nei test.
// t.Helper() fa sì che in caso di errore Go mostri la riga del test chiamante, non questa funzione.
// Usa bcrypt.MinCost per velocizzare i test (il costo non ha importanza qui).
func validToken(t *testing.T) string {
	t.Helper()
	os.Setenv("JWT_SECRET", "test-secret")
	claims := jwt.MapClaims{
		"sub":   "user-123",
		"email": "test@test.com",
		"exp":   time.Now().Add(time.Hour).Unix(), // scade tra 1 ora
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

// TestAuth_NoHeader verifica che una richiesta senza header Authorization venga bloccata con 401.
// È il caso più comune: client che dimentica di inviare il token.
func TestAuth_NoHeader(t *testing.T) {
	// wrappa un handler fittizio con il middleware — se il middleware passa la richiesta, il handler risponde 200
	h := middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil) // nessun header Authorization
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TestAuth_InvalidToken verifica che un token malformato venga rifiutato con 401.
// Copre il caso di token manomesso, troncato o generato con chiave sbagliata.
func TestAuth_InvalidToken(t *testing.T) {
	h := middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token-non-valido")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TestAuth_ValidToken verifica che un token valido venga accettato e che l'userID
// venga correttamente iniettato nel context della request.
// I handler protetti leggono l'userID dal context per sapere chi sta facendo la richiesta.
func TestAuth_ValidToken(t *testing.T) {
	token := validToken(t)

	h := middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verifica che il middleware abbia messo l'userID nel context
		userID := r.Context().Value(middleware.UserIDKey)
		if userID == nil {
			t.Error("expected userID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TestAuth_ExpiredToken verifica che un token con exp nel passato venga rifiutato.
// Il JWT contiene la scadenza al suo interno — il middleware la controlla automaticamente.
func TestAuth_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	claims := jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(), // scaduto 1 ora fa
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("test-secret"))

	h := middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
