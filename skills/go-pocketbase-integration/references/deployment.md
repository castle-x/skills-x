# Deployment & Data Management

## Table of Contents

1. [Running the Binary](#running-the-binary)
2. [Data Directory Structure](#data-directory-structure)
3. [Docker Deployment](#docker-deployment)
4. [Systemd Service](#systemd-service)
5. [Environment Variables](#environment-variables)
6. [Backup & Restore](#backup--restore)
7. [Admin UI](#admin-ui)
8. [Cross-Platform Build](#cross-platform-build)

## Running the Binary

```bash
# Default — listens on localhost:8090
./app serve

# Custom host/port
./app serve --http=0.0.0.0:8090

# With HTTPS (auto Let's Encrypt)
./app serve --https=yourdomain.com

# Show all commands
./app --help
```

PocketBase uses Cobra CLI. The `serve` command starts the HTTP server.

## Data Directory Structure

PocketBase stores all data in the `DefaultDataDir` (set in `pocketbase.Config`):

```
app_data/
├── data.db           # Main SQLite database (collections, records, auth)
├── logs.db           # Request/error logs
├── storage/          # Uploaded files (organized by collection/record)
│   └── {collection_id}/
│       └── {record_id}/
│           └── {filename}
├── types.d.ts        # Auto-generated TypeScript types for collections
└── backups/          # Admin-initiated backups
```

### Important notes

- `data.db` is the **single source of truth** — all collections, records, settings, and users
- SQLite WAL mode is used for concurrent read/write
- The data directory must be **writable** by the process
- In Docker, mount this as a **volume** to persist across container recreations

## Docker Deployment

### Multi-stage Dockerfile

```dockerfile
# Build stage
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source (includes pre-built site/dist/)
COPY . ./

RUN apk add --no-cache ca-certificates && update-ca-certificates

# Build binary
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags "-w -s" -o /app-binary ./cmd/app

# Runtime stage — scratch for minimal image
FROM scratch
COPY --from=builder /app-binary /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

VOLUME ["/app_data"]
EXPOSE 8090

ENTRYPOINT ["/app"]
CMD ["serve", "--http=0.0.0.0:8090"]
```

### Docker Compose

```yaml
services:
  app:
    build: .
    ports:
      - "8090:8090"
    volumes:
      - app_data:/app_data
    restart: unless-stopped

volumes:
  app_data:
```

### Build and run

```bash
docker build -t myapp .
docker run -p 8090:8090 -v app_data:/app_data myapp
```

## Systemd Service

```ini
[Unit]
Description=My App
After=network.target

[Service]
Type=simple
User=appuser
Group=appuser
ExecStart=/opt/app/app serve --http=0.0.0.0:8090
WorkingDirectory=/opt/app
Restart=on-failure
RestartSec=5s

# Security
LimitNOFILE=4096
PrivateTmp=true
NoNewPrivileges=true

# Environment
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

Install:

```bash
sudo cp app.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable app
sudo systemctl start app
sudo systemctl status app
```

## Environment Variables

Common patterns for PocketBase-integrated apps:

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | (empty) | Set to `dev` to enable dev mode and auto-migrations |
| `APP_URL` | (empty) | Public URL for the app (used for subpath routing) |
| `USER_EMAIL` | (empty) | Initial admin email (used in initial migration) |
| `USER_PASSWORD` | (empty) | Initial admin password |
| `DISABLE_PASSWORD_AUTH` | (empty) | Set to `true` to disable email/password login |
| `TRUSTED_AUTH_HEADER` | (empty) | HTTP header for trusted proxy authentication |

Read environment variables with optional prefix:

```go
func GetEnv(key string) (string, bool) {
    if v, ok := os.LookupEnv("MYAPP_" + key); ok {
        return v, true
    }
    return os.LookupEnv(key)
}
```

## Backup & Restore

### Via Admin UI

PocketBase admin UI at `/_/` provides backup/restore functionality.

### Manual backup

```bash
# Stop the app or use SQLite backup API
sqlite3 app_data/data.db ".backup backup.db"

# Or simply copy (safe if app is stopped)
cp -r app_data/ app_data_backup/
```

### Automated backup (via cron job in app)

```go
h.App.Cron().MustAdd("daily backup", "0 2 * * *", func() {
    // create backup using PocketBase's built-in backup
    if err := h.App.CreateBackup(context.Background(), "auto_backup"); err != nil {
        h.Logger().Error("backup failed", "err", err)
    }
})
```

### Restore

1. Stop the application
2. Replace `app_data/data.db` with the backup
3. Start the application

## Admin UI

PocketBase includes a built-in admin UI at `/_/`:

- Manage collections (schema)
- Browse/edit records
- Manage users and auth settings
- View logs
- Configure OAuth2 providers
- Create/restore backups

Access requires a **superuser** account (separate from regular users).

### Hiding admin controls

To hide admin UI controls in production:

```go
settings.Meta.HideControls = true
```

The admin UI is still accessible but some modification controls are hidden.

## Cross-Platform Build

### Makefile targets

```makefile
BINARY := app
CMD := ./cmd/app
LDFLAGS := -s -w

build-web:
	cd site && pnpm build

build: build-web
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY) $(CMD)

build-linux: build-web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY)-linux-amd64 $(CMD)

build-linux-arm64: build-web
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY)-linux-arm64 $(CMD)

build-macos: build-web
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY)-darwin-arm64 $(CMD)

build-windows: build-web
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY)-windows-amd64.exe $(CMD)

build-all: build-web
	# ... all platform combinations

dev:
	cd site && pnpm dev &
	go run -tags development $(CMD) serve

dev-backend:
	go run -tags development $(CMD) serve
```

`CGO_ENABLED=0` is critical — PocketBase uses pure-Go SQLite (`modernc.org/sqlite`), no C compiler needed.
