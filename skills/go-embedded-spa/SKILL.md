---
name: go-embedded-spa
description: This skill provides guidance for implementing Go Embedded SPA architecture - embedding React/Vue/TSX frontend static resources into Go binary using go:embed directive. Use this skill when building self-contained single-binary applications, implementing SPA with Go backend, setting up cross-platform deployable full-stack projects, or configuring static file serving with Hertz/Gin frameworks.
---

# Go Embedded SPA

## Overview

Go Embedded SPA is a technique that embeds frontend SPA (Single Page Application) static resources (React/Vue/TSX) into Go binary files using Go 1.16+ `embed` package, achieving **single-binary full-stack deployment**.

### Core Benefits

| Benefit | Description |
|---------|-------------|
| 🎯 Single File Deploy | One binary contains both frontend and backend, no nginx needed |
| 🌍 Cross-Platform | `GOOS/GOARCH` easily compiles for Linux/Mac/Windows |
| 📦 Zero Dependencies | Target machine needs no Node.js/npm |
| 🚀 Container Friendly | Dockerfile only needs `COPY + ENTRYPOINT` |
| 🔒 Resource Security | Static resources compiled into binary, tamper-proof |
| ⚡ Fast Startup | No disk I/O for loading static files |

## Architecture Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         Build Process                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   Frontend Source      Vite/Webpack        Build Output          │
│  ┌──────────────┐     ┌──────────────┐    ┌──────────────┐      │
│  │  site/src/*  │ ──► │    build     │ ──► │  site/dist/* │      │
│  │  (React TSX) │     │              │    │ (static files)│      │
│  └──────────────┘     └──────────────┘    └──────────────┘      │
│                                                  │                │
│                                                  ▼                │
│   Go Source            go build           Final Binary           │
│  ┌──────────────┐     ┌──────────────┐    ┌──────────────┐      │
│  │    *.go      │ ──► │   compile    │ ◄── │ //go:embed   │      │
│  │  embed.go    │     │              │    │   all:dist   │      │
│  └──────────────┘     └──────────────┘    └──────────────┘      │
│                              │                                    │
│                              ▼                                    │
│                       ┌──────────────┐                           │
│                       │ bin/app      │ (contains frontend)       │
│                       └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

## Implementation Guide

### Step 1: Project Structure

Create the following directory structure:

```
project/
├── site/                    # Frontend project
│   ├── src/                 # Source code
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   └── ...
│   ├── dist/                # Build output (embedded)
│   ├── embed.go             # Go embed directive
│   ├── package.json
│   ├── vite.config.ts
│   └── index.html
├── pkg/
│   └── siteserver/          # Static file server
│       └── siteserver.go
├── internal/
│   └── app/
│       └── app.go           # Application integration
├── Makefile
└── go.mod
```

### Step 2: Create embed.go

Create `site/embed.go` with the embed directive:

```go
package site

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distDir embed.FS

// DistDirFS returns the embedded frontend resource filesystem
// Usage: siteserver.StaticFS(h, site.DistDirFS)
var DistDirFS, _ = fs.Sub(distDir, "dist")
```

**Key Points:**
- `//go:embed all:dist` - The `all:` prefix embeds ALL files including those starting with `.` or `_`
- `embed.FS` - Go 1.16+ embedded read-only filesystem type
- `fs.Sub(distDir, "dist")` - Creates sub-filesystem, removes `dist/` prefix for clean access

### Step 3: Create Static File Server

Create `pkg/siteserver/siteserver.go` for serving embedded files. See `references/siteserver.md` for the complete implementation.

**Core Logic:**
1. Pre-load `index.html` for SPA route fallback
2. Create standard library file server from embed.FS
3. Register NoRoute handler (catches all unmatched routes)
4. Static resource detection (path contains `.`)
5. SPA fallback (return index.html for non-static routes)

### Step 4: Application Integration

In the main application, register static file server AFTER API routes:

```go
import (
    "your-project/pkg/siteserver"
    "your-project/site"
)

func NewApp() {
    h := server.Default()
    
    // 1. Register API routes FIRST
    h.POST("/apis/v1/users", userHandler)
    h.GET("/apis/v1/data", dataHandler)
    
    // 2. Register static file server LAST (as fallback)
    if err := siteserver.StaticFS(h, site.DistDirFS); err != nil {
        log.Warn("Failed to register static file server: %v", err)
    }
    
    h.Spin()
}
```

**Order is critical:** API routes must be registered before static file server to ensure API matching priority.

### Step 5: Build Configuration

#### Makefile

```makefile
build-web:  ## Build frontend
	cd site && npm run build

build-backend:  ## Build backend
	go build -o bin/app ./cmd/app

build: build-web build-backend  ## Build all (frontend first!)

clean:  ## Clean build artifacts
	rm -rf site/dist bin/app
```

**Build order MUST be:** `build-web` → `build-backend`

#### Vite Configuration (site/vite.config.ts)

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',              // Output dir matches embed directive
    emptyDirBeforeWrite: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/apis': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
})
```

### Step 6: Cross-Platform Build

```bash
# Build for multiple platforms
GOOS=linux   GOARCH=amd64 go build -o bin/app-linux ./cmd/app
GOOS=darwin  GOARCH=arm64 go build -o bin/app-macos ./cmd/app
GOOS=windows GOARCH=amd64 go build -o bin/app.exe ./cmd/app
```

### Step 7: Container Deployment

Minimal Dockerfile:

```dockerfile
FROM scratch
COPY app /app
ENTRYPOINT ["/app"]
```

Or with Alpine base:

```dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY app /app
ENTRYPOINT ["/app"]
```

## Request Handling Flow

```
Browser Request
      │
      ▼
┌─────────────────────────────────────┐
│           HTTP Server               │
├─────────────────────────────────────┤
│  Route Matching                     │
│  ├── /apis/*  → API Handler         │
│  ├── /ws      → WebSocket Handler   │
│  └── Others   → NoRoute (siteserver)│
└─────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────┐
│       NoRoute Handler Logic         │
├─────────────────────────────────────┤
│  GET /assets/index-abc123.js        │
│  → Has extension → embed.FS read    │
│  → Return JS file                   │
│                                     │
│  GET /dashboard/settings            │
│  → No extension → SPA fallback      │
│  → Return index.html                │
└─────────────────────────────────────┘
```

## Caching Strategy

| Path Pattern | Cache-Control | Reason |
|--------------|---------------|--------|
| `/assets/*` | `public, max-age=31536000, immutable` | Files have hash in name, safe for long cache |
| `/index.html` | `no-cache` | Entry file must always be fresh |
| Other static | Default | Standard browser caching |

## Troubleshooting

### Common Issues

1. **Empty dist directory error**
   - Ensure `make build-web` runs BEFORE `make build-backend`
   - Check that `site/dist/` exists and contains files

2. **Static files not found**
   - Verify `//go:embed all:dist` path is correct relative to embed.go location
   - Check `fs.Sub()` prefix matches actual directory structure

3. **API routes not matching**
   - Ensure API routes are registered BEFORE `siteserver.StaticFS()`
   - Check route patterns don't conflict

4. **SPA routes return 404**
   - Verify NoRoute handler returns index.html for non-file paths
   - Check path detection logic (looking for `.` in filename)

## Resources

### references/
- `siteserver.md` - Complete static file server implementation for Hertz framework

### assets/
- `embed.go.tmpl` - Template for embed.go file
- `siteserver.go.tmpl` - Template for static file server
- `vite.config.ts.tmpl` - Template for Vite configuration
