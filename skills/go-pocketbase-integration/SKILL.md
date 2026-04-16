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
    "path/filepath"

    "your-project/internal/server"
    _ "your-project/internal/migrations"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func getDataDir() string {
    if dir := os.Getenv("DATA_DIR"); dir != "" {
        return dir
    }
    home, err := os.UserHomeDir()
    if err != nil {
        return "app_data"
    }
    return filepath.Join(home, ".myapp")
}

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

### 1. Password validation: use ValidatePassword, NOT bcrypt

```go
// WRONG — record.GetString("password") returns empty string for hash fields
hash := record.GetString("password")
bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) // always fails

// CORRECT — use PocketBase's built-in method
record.ValidatePassword("plaintext") // true/false
```

### 2. Default password minimum length is 8 characters

PocketBase enforces a minimum password length of **8 characters** for auth collections. If `gopb.Options.DefaultPassword` is shorter than 8 characters, `s.EnsureDefaults()` will fail with an error like:

```
ERROR gopb: failed to create default superuser
└─ {"error":{"password":"Must be at least 8 character(s)"}}
```

```go
// WRONG — too short, will fail silently or log errors
srv := gopb.New(app, gopb.Options{
    DefaultEmail:    "admin@myapp.local",
    DefaultPassword: "myapp123",  // only 7 chars!
})

// CORRECT — minimum 8 characters
srv := gopb.New(app, gopb.Options{
    DefaultEmail:    "admin@myapp.local",
    DefaultPassword: "myapp1234", // 8+ chars
})
```

### 3. Never register routes that conflict with PocketBase built-ins

PocketBase registers its own routes (e.g., `GET /api/health`, `POST /api/collections/:collection/auth-with-password`). Adding a custom route with the same method+path causes a runtime panic:

```
panic: pattern "GET /api/health" conflicts with pattern "GET /api/health"
```

**Rule**: Only register custom routes under paths that don't overlap with PocketBase's API. Common safe prefixes:
- `/api/v1/` or `/api/setup/` (used by gopb)
- `/api/{your-business-domain}/`

```go
// WRONG — conflicts with PocketBase's built-in health endpoint
e.Router.GET("/api/health", handler) // PANIC at runtime

// CORRECT — use a unique prefix
apiAuth := e.Router.Group("/api")
apiAuth.Bind(apis.RequireAuth())
apiAuth.GET("/items", s.handleListItems)
```

### 4. Two collections: `_superusers` vs `users`

- `_superusers` is for Admin UI (`/_/`) access only. It does NOT support public API authentication.
- `users` is for frontend app login via `pb.collection("users").authWithPassword()`.
- Always create **both** a superuser and a user with the same credentials for first-run setup.
- The `handleChangePassword` endpoint syncs passwords across both collections.

### 5. Login state is NOT shared

The `_superusers` and `users` collections have independent auth sessions. Logging into the app does NOT log you into the Admin UI and vice versa. This is by design in PocketBase. Mitigation: keep passwords in sync (which `gopb` does automatically).

### 6. e.Auth in protected routes

In routes protected by `apis.RequireAuth()`, `e.Auth` is the authenticated `*core.Record`. Access `e.Auth.Id`, `e.Auth.Email()`, `e.Auth.GetString("field")` directly.

### 7. Data directory should be configurable, not hardcoded

Avoid hardcoding `DefaultDataDir` to `"app_data"`. Instead, use `~/.appname/` as default with environment variable override:

```go
func getDataDir() string {
    if dir := os.Getenv("DATA_DIR"); dir != "" {
        return dir
    }
    home, err := os.UserHomeDir()
    if err != nil {
        return "app_data" // fallback
    }
    return filepath.Join(home, ".appname")
}

func main() {
    app := pocketbase.NewWithConfig(pocketbase.Config{
        DefaultDataDir: getDataDir(),
        DefaultDev:     isDev,
    })
    // ...
}
```

This ensures:
- Dev data doesn't pollute the project directory
- Production data lives in a predictable, XDG-compliant location
- `make clean` doesn't accidentally delete production data

### 8. Vite dev server must be configured for external access

When accessing the frontend through a tunnel or reverse proxy (e.g., `https://yourdomain.top`), Vite will block requests with: `"This host isyourdomain.top" is not allowed`.

Three settings must be configured in `vite.config.ts`:

```typescript
server: {
    host: "0.0.0.0",              // 1. Listen on all interfaces, not just localhost
    port: 3000,
    strictPort: true,
    allowedHosts: ["yourdomain.top"], // 2. Allow your tunnel/proxy domain
    proxy: {
        "/api": { target: backendTarget, changeOrigin: true },
        "/_/":  { target: backendTarget, changeOrigin: true },
    },
}
```

For the PocketBase backend, also proxy the `/_/` route in both dev and production for Admin UI access.

### 9. Makefile dev target must run backend and frontend in parallel

A common mistake is `dev: dev-backend dev-frontend` which runs them sequentially (the first blocks forever). Use background jobs:

```makefile
dev-backend:
    ENV=dev go run -tags development ./cmd/server serve --http=localhost:8180

dev-frontend:
    cd site && npx vite

dev:
    @echo "Starting backend and frontend in parallel..."
    @$(MAKE) dev-backend & $(MAKE) dev-frontend & wait
```

### 10. `EnsureDefaults()` placement in OnServe hook

Call `s.EnsureDefaults()` **after** route registration but **before** `s.serveFrontend(e)`. If called too early (before PocketBase finishes its own initialization), it may fail. The correct order:

```go
s.OnServe().BindFunc(func(e *core.ServeEvent) error {
    s.RegisterSetupRoutes(e)     // 1. Setup routes first
    // ... custom routes ...
    s.EnsureDefaults()            // 2. Create default users
    s.serveFrontend(e)            // 3. Serve frontend last (catch-all route)
    return e.Next()
})
```

### 11. Frontend authentication flow: check setup status after login

After successful login, always check whether the user needs to change their default password before granting access to the dashboard:

```typescript
const authData = await pb.collection("users").authWithPassword(email, password);
setAuth({ id: authData.record.id, email: authData.record.email, ...authData.record }, authData.token);

// Check if first-run password change is required
const status = await checkSetupStatus();
if (status.needsPasswordChange) {
    navigate("/change-password", { replace: true });
} else {
    navigate("/", { replace: true });
}
```

### 12. Auth guard must sync PocketBase authStore events to Zustand

The `AuthGuard` component must listen to `pb.authStore.onChange` and sync logout to Zustand. Otherwise, 401 responses from the API won't trigger a UI redirect:

```typescript
useEffect(() => {
    const unsubscribe = pb.authStore.onChange(() => {
        if (!pb.authStore.isValid && isAuthenticated) {
            logout();
            navigate("/login", { replace: true });
        }
    });
    return () => unsubscribe();
}, [isAuthenticated, logout, navigate]);
```

### 13. Auth state restoration must happen before React renders

The `pb-client.ts` module restores auth from localStorage on import (before any React component mounts). The `AUTH_STORAGE_KEY` in `pb-client.ts` **must match** the `name` in Zustand's `persist()` middleware:

```typescript
// pb-client.ts
const AUTH_STORAGE_KEY = "auth-storage";

// useAuth.ts
persist(/* ... */, { name: "auth-storage" }) // must match!
```

## Frontend Integration

See `references/frontend-patterns.md` for PocketBase JS SDK setup, auth state management, and SPA embedding patterns.

## References

- `references/project-structure.md` — Project layout and file organization
- `references/backend-patterns.md` — Routes, middleware, DB, migrations, hooks, cron, build tags
- `references/frontend-patterns.md` — SPA embedding, PocketBase JS SDK, auth patterns
- `references/deployment.md` — Docker, systemd, data management, backup
