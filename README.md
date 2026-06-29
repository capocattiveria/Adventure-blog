# Adventure Blog

Web application per condividere avventure di viaggio. Backend in Go, frontend mobile in React Native/Expo, database PostgreSQL e cache Redis.

---

## Stack

| Layer | Tecnologia |
|---|---|
| Backend API | Go 1.22 |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Hot reload | Air |
| Migrazioni DB | golang-migrate |
| Mobile | React Native + Expo |

---

## Prerequisiti

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installato e in esecuzione
- [VS Code](https://code.visualstudio.com/) con l'estensione [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

Non serve installare Go, Node o PostgreSQL in locale — tutto gira nel container.

---

## Avvio rapido

### 1. Apri il progetto in VS Code

```bash
git clone <url-repo>
cd adventure-blog
code .
```

### 2. Avvia il Dev Container

Quando VS Code rileva il file `.devcontainer/devcontainer.json`, mostra una notifica in basso a destra:

> **Reopen in Container**

Clicca il pulsante. In alternativa: `Ctrl+Shift+P` → `Dev Containers: Reopen in Container`.

VS Code costruirà il container, avvierà PostgreSQL e Redis, e installerà tutti gli strumenti Go automaticamente. La prima volta richiede qualche minuto. Dai rebuild successivi sarà molto più veloce grazie ai volumi Docker che fanno da cache.

### 3. Avvia il backend

Apri un terminale dentro VS Code (`Ctrl+` `` ` ``) — sei già dentro il container.

```bash
cd /workspace/adventure-blog/backend
air
```

`air` compila e avvia il server Go, e lo riavvia automaticamente ad ogni salvataggio. L'API sarà disponibile su `http://localhost:8080`.

---

## Struttura del progetto

```
adventure-blog/
├── backend/              # API Go
│   ├── Dockerfile.dev    # Immagine Docker per sviluppo
│   └── ...
├── mobile/               # App React Native (Expo)
│   └── ...
├── .devcontainer/
│   ├── devcontainer.json # Configurazione Dev Container
│   └── post-create.sh    # Script di setup (Go tools, dipendenze)
└── docker-compose.yml    # Orchestrazione container
```

---

## Servizi disponibili

| Servizio | URL / Host | Note |
|---|---|---|
| Go API | `http://localhost:8080` | avviato manualmente con `air` |
| PostgreSQL | `localhost:5432` | avviato automaticamente |
| Redis | `localhost:6379` | avviato automaticamente |
| Expo Web | `http://localhost:8081` | avviato manualmente con `npx expo start` |

### Credenziali PostgreSQL (default)

```
host:     localhost
port:     5432
database: adventure_blog
user:     adventure
password: adventure_secret
```

Personalizzabili tramite variabili d'ambiente nel file `.env` (crea il file copiando i valori da `docker-compose.yml`).

---

## Migrazioni database

Le migrazioni si trovano in `backend/migrations/`. Per applicarle:

```bash
# Applica tutte le migrazioni pendenti
migrate -path backend/migrations -database "$DATABASE_URL" up

# Torna indietro di una migrazione
migrate -path backend/migrations -database "$DATABASE_URL" down 1
```

La variabile `$DATABASE_URL` è già disponibile nel container:
```
postgres://adventure:adventure_secret@postgres:5432/adventure_blog?sslmode=disable
```

---

## Comandi utili

```bash
# Controlla i log di PostgreSQL
docker compose logs postgres

# Controlla i log di Redis
docker compose logs redis

# Accedi al database con psql
psql "$DATABASE_URL"

# Formatta il codice Go
goimports -w .

# Esegui il linter
golangci-lint run ./...

# Build manuale (senza hot reload)
go build ./...

# Esegui i test
go test ./...
```

---

## Rebuild del Dev Container

Se aggiungi nuovi strumenti o modifichi `devcontainer.json`:

`Ctrl+Shift+P` → `Dev Containers: Rebuild Container`

I moduli Go e i binari degli strumenti sono salvati in volumi Docker persistenti, quindi il rebuild non riscarica tutto da zero.
