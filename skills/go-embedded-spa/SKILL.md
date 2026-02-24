---
name: go-embedded-spa
description: This skill provides guidance for implementing Go Embedded SPA architecture - embedding React/Vue/TSX frontend static resources into Go binary using go:embed directive. Use this skill when building self-contained single-binary applications, implementing SPA with Go backend, setting up cross-platform deployable full-stack projects, or configuring static file serving with Go standard library net/http.
---

# Go Embedded SPA

## Overview

Go Embedded SPA is a technique that embeds frontend SPA (Single Page Application) static resources (React/Vue/TSX) into Go binary files using Go 1.16+ `embed` package, achieving **single-binary full-stack deployment**.

### Core Benefits

| Benefit | Description |
|---------|-------------|
| 🎯 Single File Deploy | One binary contains both frontend and backend, no nginx needed |
| 🌍 Cross-Platform | `GOOS/GOARCH` easily compiles for Linux/Mac/Windows |
| 📦 Zero Dependencies | Target machine needs no Node.js/npm, uses Go standard library only |
| 🚀 Container Friendly | Dockerfile only needs `COPY + ENTRYPOINT` |
| 🔒 Resource Security | Static resources compiled into binary, tamper-proof |
| ⚡ Fast Startup | No disk I/O for loading static files |

## Project Structure

```
project/
├── go.mod
├── Makefile
├── site/                    # Frontend project
│   ├── embed.go             # Go embed directive
│   ├── src/                 # React/Vue source
│   ├── dist/                # Build output (embedded)
│   ├── package.json
│   ├── vite.config.ts
│   └── index.html
├── pkg/
│   └── siteserver/          # Static file server
│       └── siteserver.go
└── cmd/
    └── app/
        └── main.go
```

## Implementation Steps

### Step 1: Create embed.go

Create `site/embed.go` to declare embed directive. See `references/embed.md` for complete code.

Key points:
- Use `//go:embed all:dist` to embed all files including hidden files
- Use `fs.Sub()` to remove `dist/` prefix

### Step 2: Create Static File Server

Create `pkg/siteserver/siteserver.go`. See `references/siteserver.md` for complete implementation using Go standard `net/http`.

Core logic:
1. Pre-load `index.html` for SPA fallback
2. Create `http.FileServer` from embed.FS
3. Detect static resources by file extension
4. Return `index.html` for SPA routes (no file extension)

### Step 3: Application Integration

```go
package main

import (
    "log"
    "net/http"
    
    "your-project/pkg/siteserver"
    "your-project/site"
)

func main() {
    mux := http.NewServeMux()
    
    // 1. Register API routes FIRST
    mux.HandleFunc("/apis/v1/health", healthHandler)
    mux.HandleFunc("/apis/v1/data", dataHandler)
    
    // 2. Wrap with static file server (as fallback)
    handler := siteserver.WrapHandler(mux, site.DistDirFS)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

**Order is critical:** API routes must be registered before static file server.

### Step 4: Development

**Recommended:** Start both backend and frontend with one command:

```bash
make dev
# Backend:  http://localhost:8080 (Go)
# Frontend: http://localhost:5173 (Vite with hot reload)
# API Proxy: /api/* -> localhost:8080
```

Press `Ctrl+C` to stop both servers.

### Step 5: Build for Production

Build order: **frontend first, then backend**

```bash
make build          # Build both (frontend + backend)
# Or separately:
make build-web      # npm run build → site/dist/
make build-backend  # go build (embeds dist/)
```

### Step 6: Cross-Platform Build

```bash
make build-linux        # Linux amd64
make build-linux-arm64  # Linux arm64
make build-macos        # macOS Intel
make build-macos-arm64  # macOS Apple Silicon
make build-windows      # Windows amd64
make build-all          # All platforms
```

## Request Handling Flow

```
Browser Request
      │
      ▼
┌─────────────────────────────────────┐
│  http.ServeMux Route Matching       │
│  ├── /apis/*  → API Handler         │
│  └── Others   → Static Handler      │
└─────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────┐
│  Static Handler                     │
│  ├── Has extension → serve file     │
│  └── No extension → return index.html│
└─────────────────────────────────────┘
```

## Caching Strategy

| Path | Cache-Control | Reason |
|------|---------------|--------|
| `/assets/*` | `max-age=31536000, immutable` | Files have hash in name |
| `/index.html` | `no-cache` | Entry must be fresh |

## Container Deployment

Minimal Dockerfile:

```dockerfile
FROM scratch
COPY app /app
ENTRYPOINT ["/app"]
```

## Troubleshooting

1. **Empty dist error** → Run `make build-web` before `make build-backend`
2. **Static files 404** → Check `//go:embed all:dist` path relative to embed.go
3. **API not matching** → Register API routes BEFORE wrapping with siteserver
4. **SPA routes 404** → Verify handler returns index.html for non-file paths

## References

- `go-dependencies.md` - Go module dependencies (standard library only)
- `embed.md` - Complete embed.go implementation
- `siteserver.md` - Static file server using Go standard net/http
- `vite-config.md` - Vite configuration for development proxy
- `makefile.md` - Complete Makefile with cross-platform build
