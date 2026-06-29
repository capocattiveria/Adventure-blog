#!/usr/bin/env bash
set -e

echo "──────────────────────────────────────────"
echo " Adventure Blog — post-create setup"
echo "──────────────────────────────────────────"

# Docker named volumes are initialised as root; fix ownership so the vscode
# user can write to the Go cache and module directories.
sudo chown -R vscode:vscode /home/vscode/.cache/go-build 2>/dev/null || true
sudo chown -R vscode:vscode /home/vscode/go 2>/dev/null || true

# ── Go tools ───────────────────────────────────────────────────────────────
# Installa solo se il binario non esiste già (il volume go-bin persiste tra rebuild)
echo "→ Installing Go tools (skipping if already cached)..."
command -v air        >/dev/null || go install github.com/air-verse/air@latest
command -v goimports  >/dev/null || go install golang.org/x/tools/cmd/goimports@latest
command -v golangci-lint >/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
command -v migrate    >/dev/null || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ── Backend deps ───────────────────────────────────────────────────────────
if [ -f backend/go.mod ]; then
  echo "→ Downloading Go modules (skipping if already cached)..."
  # go mod download è no-op se i moduli sono già nel volume go-pkg-mod
  cd backend && go mod download && cd ..
fi

# ── Node / Expo ────────────────────────────────────────────────────────────
echo "→ Installing Expo CLI globally..."
npm install -g expo-cli @expo/ngrok

if [ -f mobile/package.json ]; then
  echo "→ Installing mobile dependencies..."
  cd mobile && npm install && cd ..
fi

echo "✓ Setup complete. Happy coding!"