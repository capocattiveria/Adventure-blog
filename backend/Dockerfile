# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /api ./cmd/api

# Include migrate binary for running DB migrations inside the container
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S adventure && adduser -S adventure -G adventure

COPY --from=builder /api /usr/local/bin/api
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY migrations /migrations

USER adventure
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api"]
