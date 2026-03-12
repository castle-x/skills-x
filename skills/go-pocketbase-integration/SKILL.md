---
name: go-pocketbase-integration
description: Integrate PocketBase as a Go library using the github.com/castle-x/goutils/pocketbase (gopb) package to build single-binary full-stack applications. Use when building Go applications that need user authentication, embedding PocketBase into Go binary, registering custom API routes, managing default users, serving embedded SPA frontend, or deploying single-binary applications. NOT for using PocketBase as a standalone separate process.
---

# Go PocketBase Integration

Embed PocketBase as a Go library using the **`gopb`** package (`github.com/castle-x/goutils/pocketbase`) to produce a **single-binary full-stack application** with built-in auth, SQLite database, admin UI, and custom API routes.

## Architecture Overview

```
Single Binary
├── gopb.AppServer (wraps PocketBase core.App)
│   ├── Default user initialization (superuser + app user)
│   ├── Setup routes (status check + change password)
│   └── SPA serving helpers (production + dev proxy)
├── Custom API Routes (business logic)
├── Migrations (schema version control)
└── Embedded SPA Frontend (go:embed)
```

PocketBase runs **in-process** — no separate service, no HTTP calls for auth validation.

## Quick Start

### 1. Add dependencies

```bash
go get github.com/castle-x/goutils/pocketbase@latest
go get github.com/pocketbase/pocketbase@latest
```

### 2. Entry point

```go
package main

import (
    "log"
    "os"

    "your-project/internal/server"
    _ "your-project/internal/migrations"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
    isDev := os.Getenv("ENV") == "dev"

    app := pocketbase.NewWithConfig(pocketbase.Config{
        DefaultDataDir: getDataDir(),
        DefaultDev:     isDev,
    })

    migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
        Automigrate: isDev,
        Dir:         "internal/migrations",
    })

    srv := server.New(app)
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### 3. AppServer — use gopb building blocks

```go
package server

import (
    gopb "github.com/castle-x/goutils/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)

type AppServer struct {
    *gopb.AppServer
    dataPath string
}

func New(app core.App) *AppServer {
    srv := gopb.New(app, gopb.Options{
        DefaultEmail:    "admin@myapp.local",
        DefaultPassword: "myapp123",
    })
    return &AppServer{AppServer: srv}
}

func (s *AppServer) Start() error {
    s.OnServe().BindFunc(func(e *core.ServeEvent) error {
        // 1. Setup routes (status check + change password)
        s.RegisterSetupRoutes(e)

        // 2. Business routes
        api := e.Router.Group("/api")
        api.Bind(apis.RequireAuth())
        api.GET("/items", s.handleListItems)

        // 3. Create default users if first run
        s.EnsureDefaults()

        // 4. Serve frontend (use build tags to switch)
        s.serveFrontend(e)

        return e.Next()
    })
    return s.AppServer.Start()
}
```

### 4. SPA serving (build tags)

```go
//go:build !development
// file: server_production.go

package server

import (
    gopb "github.com/castle-x/goutils/pocketbase"
    "github.com/pocketbase/pocketbase/core"
    "your-project/site"
)

func (s *AppServer) serveFrontend(se *core.ServeEvent) {
    gopb.ServeSPA(se, site.DistDirFS, []string{"/assets/", "/static/"})
}
```

```go
//go:build development
// file: server_development.go

package server

import (
    gopb "github.com/castle-x/goutils/pocketbase"
    "github.com/pocketbase/pocketbase/core"
)

func (s *AppServer) serveFrontend(se *core.ServeEvent) {
    gopb.ServeDevProxy(se, "localhost:5173")
}
```

## gopb API Reference

### Core

| Function | Description |
|----------|-------------|
| `gopb.New(app, opts...)` | Create AppServer wrapping core.App |
| `s.Start()` | Launch PocketBase (blocks) |
| `s.Opts()` | Get resolved Options |

### Default User Management

| Function | Description |
|----------|-------------|
| `s.EnsureDefaults()` | Create default superuser + user if collections empty |
| `s.IsDefaultPassword(record)` | Check if record uses default password |

### Setup Routes

| Function | Description |
|----------|-------------|
| `s.RegisterSetupRoutes(se)` | Register GET /status + POST /change-password |

Endpoints (under `Options.SetupRoutePrefix`, default `/api/setup`):
- `GET /status` → `{"needsPasswordChange": bool}`
- `POST /change-password` → accepts `{"password", "passwordConfirm"}`

The change-password endpoint **syncs passwords to `_superusers`** — any superuser still using the default password gets updated, keeping Admin UI access in sync.

### SPA Helpers

| Function | Description |
|----------|-------------|
| `gopb.ServeSPA(se, distFS, staticPaths)` | Serve embedded SPA with cache + fallback |
| `gopb.ServeDevProxy(se, host)` | Proxy to Vite dev server |

## Key Pitfalls & Resolved Issues

### Password validation: use ValidatePassword, NOT bcrypt

```go
// WRONG — record.GetString("password") returns empty string for hash fields
hash := record.GetString("password")
bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) // always fails

// CORRECT — use PocketBase's built-in method
record.ValidatePassword("plaintext") // true/false
```

### Two collections: `_superusers` vs `users`

- `_superusers` is for Admin UI (`/_/`) access only. It does NOT support public API authentication.
- `users` is for frontend app login via `pb.collection("users").authWithPassword()`.
- Always create **both** a superuser and a user with the same credentials for first-run setup.
- The `handleChangePassword` endpoint syncs passwords across both collections.

### Login state is NOT shared

The `_superusers` and `users` collections have independent auth sessions. Logging into the app does NOT log you into the Admin UI and vice versa. This is by design in PocketBase. Mitigation: keep passwords in sync (which `gopb` does automatically).

### e.Auth in protected routes

In routes protected by `apis.RequireAuth()`, `e.Auth` is the authenticated `*core.Record`. Access `e.Auth.Id`, `e.Auth.Email()`, `e.Auth.GetString("field")` directly.

## Frontend Integration

See `references/frontend-patterns.md` for PocketBase JS SDK setup, auth state management, and SPA embedding patterns.

## References

- `references/project-structure.md` — Project layout and file organization
- `references/backend-patterns.md` — Routes, middleware, DB, migrations, hooks, cron, build tags
- `references/frontend-patterns.md` — SPA embedding, PocketBase JS SDK, auth patterns
- `references/deployment.md` — Docker, systemd, data management, backup
