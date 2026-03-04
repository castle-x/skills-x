# Backend Patterns

## Table of Contents

1. [Entry Point](#entry-point)
2. [AppServer Pattern (with gopb)](#appserver-pattern-with-gopb)
3. [Route Registration](#route-registration)
4. [Authentication & Middleware](#authentication--middleware)
5. [Database Operations](#database-operations)
6. [Migrations](#migrations)
7. [Record Hooks](#record-hooks)
8. [Cron Jobs](#cron-jobs)
9. [Build Tags](#build-tags)
10. [Custom CLI Commands](#custom-cli-commands)
11. [Settings Configuration](#settings-configuration)
12. [Pitfalls & Resolved Issues](#pitfalls--resolved-issues)

## Entry Point

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
        DefaultDataDir: "app_data",
        DefaultDev:     isDev,
    })

    // auto-create migration files when changing collections in admin UI
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

## AppServer Pattern (with gopb)

Use `gopb.AppServer` from `github.com/castle-x/go-pocketbase` as the base. It wraps `core.App` and provides default user management, setup routes, and SPA helpers.

```go
package server

import (
    gopb "github.com/castle-x/go-pocketbase"
    "github.com/pocketbase/pocketbase/core"
)

type AppServer struct {
    *gopb.AppServer
    dataPath string // example: app-specific field
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
        // 1. Setup routes (default password check + change-password)
        s.RegisterSetupRoutes(e)
        // 2. Register business routes
        s.registerRoutes(e)
        // 3. Register cron jobs
        s.registerCronJobs()
        // 4. Create default users if first run
        s.EnsureDefaults()
        // 5. Serve frontend
        s.serveFrontend(e)
        return e.Next()
    })

    // Register record hooks (outside OnServe)
    s.registerHooks()

    return s.AppServer.Start()
}
```

**gopb provides:**
- `s.RegisterSetupRoutes(e)` — registers `/api/setup/status` and `/api/setup/change-password`
- `s.EnsureDefaults()` — creates default superuser + user if collections empty
- `s.IsDefaultPassword(record)` — checks if record uses the default password
- `gopb.ServeSPA(e, distFS, staticPaths)` — serves embedded SPA in production
- `gopb.ServeDevProxy(e, host)` — proxies to Vite in development

## Route Registration

```go
func (s *AppServer) registerRoutes(se *core.ServeEvent) error {
    // Auth-required routes
    apiAuth := se.Router.Group("/api/v1")
    apiAuth.Bind(apis.RequireAuth())

    apiAuth.GET("/items", h.listItems)
    apiAuth.POST("/items", h.createItem)
    apiAuth.GET("/items/{id}", h.getItem)
    apiAuth.PUT("/items/{id}", h.updateItem)
    apiAuth.DELETE("/items/{id}", h.deleteItem)

    // Public routes (no auth)
    apiPublic := se.Router.Group("/api/v1")
    apiPublic.GET("/health", func(e *core.RequestEvent) error {
        return e.JSON(200, map[string]string{"status": "ok"})
    })

    return nil
}

// Handler example — access auth user and query params
func (s *AppServer) listItems(e *core.RequestEvent) error {
    userID := e.Auth.Id
    query := e.Request.URL.Query()
    page := query.Get("page")

    records, err := s.App.FindAllRecords("items",
        // PocketBase filter expressions
    )
    if err != nil {
        return e.JSON(500, map[string]string{"error": err.Error()})
    }
    return e.JSON(200, records)
}

// Handler with request body
func (s *AppServer) createItem(e *core.RequestEvent) error {
    data := struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }{}
    if err := e.BindBody(&data); err != nil {
        return e.JSON(400, map[string]string{"error": err.Error()})
    }

    col, err := s.App.FindCollectionByNameOrId("items")
    if err != nil {
        return e.JSON(500, map[string]string{"error": err.Error()})
    }

    rec := core.NewRecord(col)
    rec.Set("title", data.Title)
    rec.Set("content", data.Content)
    rec.Set("user", e.Auth.Id)

    if err := s.App.Save(rec); err != nil {
        return e.JSON(500, map[string]string{"error": err.Error()})
    }
    return e.JSON(200, rec)
}
```

## Authentication & Middleware

### Built-in auth (replaces custom middleware)

```go
// Require authentication for a route group
authGroup := se.Router.Group("/api/protected")
authGroup.Bind(apis.RequireAuth())

// Inside handler — e.Auth is the authenticated record
func handler(e *core.RequestEvent) error {
    user := e.Auth          // *core.Record
    userID := e.Auth.Id     // string
    email := e.Auth.Email() // string
    role := e.Auth.GetString("role")
    return e.Next()
}
```

### Custom global middleware

```go
func (s *AppServer) registerMiddlewares(se *core.ServeEvent) {
    // Example: trusted header auth (reverse proxy scenario)
    if trustedHeader := os.Getenv("TRUSTED_AUTH_HEADER"); trustedHeader != "" {
        se.Router.BindFunc(func(e *core.RequestEvent) error {
            if e.Auth != nil {
                return e.Next()
            }
            email := e.Request.Header.Get(trustedHeader)
            if email != "" {
                e.Auth, _ = e.App.FindFirstRecordByData("users", "email", email)
            }
            return e.Next()
        })
    }
}
```

## Database Operations

All operations are in-process — no HTTP calls.

```go
// Find all records in a collection
records, err := s.App.FindAllRecords("collection_name")

// Find by ID
record, err := s.App.FindRecordById("collection_name", "record_id")

// Find first record matching filter (uses PocketBase filter syntax)
record, err := s.App.FindFirstRecordByFilter("users",
    "email = {:email}",
    dbx.Params{"email": "user@example.com"},
)

// Find by data field
record, err := s.App.FindFirstRecordByData("users", "email", "user@example.com")

// Count records
total, err := s.App.CountRecords("users")

// Create record
col, _ := s.App.FindCollectionByNameOrId("items")
rec := core.NewRecord(col)
rec.Set("title", "Hello")
rec.Set("user", userID)
s.App.Save(rec)

// Update record
rec.Set("title", "Updated")
s.App.Save(rec)

// Delete record
s.App.Delete(rec)

// Raw SQL via dbx
var result struct {
    Count int `db:"count"`
}
s.App.DB().NewQuery("SELECT COUNT(*) as count FROM items WHERE user = {:uid}").
    Bind(dbx.Params{"uid": userID}).
    One(&result)

// Find collection (cached, faster for repeated lookups)
col, err := s.App.FindCachedCollectionByNameOrId("items")
```

Import `github.com/pocketbase/dbx` for `dbx.Params`.

## Migrations

### Structure

```go
package migrations

import (
    "github.com/pocketbase/pocketbase/core"
    m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
    m.Register(func(app core.App) error {
        // UP migration

        // Configure settings
        settings := app.Settings()
        settings.Meta.AppName = "My App"
        settings.Meta.HideControls = true
        if err := app.Save(settings); err != nil {
            return err
        }

        // Create superuser
        superuserCol, _ := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
        su := core.NewRecord(superuserCol)
        su.SetEmail("admin@example.com")
        su.SetRandomPassword()
        return app.Save(su)
    }, func(app core.App) error {
        // DOWN migration (optional, can be nil)
        return nil
    })
}
```

### Collection schema via JSON snapshot

For complex schemas, use JSON snapshot migrations (auto-generated by `migratecmd`):

```go
func init() {
    m.Register(func(app core.App) error {
        jsonData := `[{"id":"...","name":"items","type":"base","fields":[...]}]`
        return app.ImportCollectionsByMarshaledJSON([]byte(jsonData), false)
    }, nil)
}
```

### Auth collection settings

```go
func setAuthSettings(app core.App) error {
    usersCol, _ := app.FindCollectionByNameOrId("users")
    usersCol.PasswordAuth.Enabled = true
    usersCol.PasswordAuth.IdentityFields = []string{"email"}
    // Set access rules
    rule := "@request.auth.id != \"\""
    usersCol.ListRule = &rule
    usersCol.ViewRule = &rule
    return app.Save(usersCol)
}
```

## Record Hooks

```go
func (s *AppServer) registerHooks() {
    // Before create
    s.App.OnRecordCreate("users").BindFunc(func(e *core.RecordEvent) error {
        if e.Record.GetString("role") == "" {
            e.Record.Set("role", "user")
        }
        return e.Next()
    })

    // After create
    s.App.OnRecordAfterCreateSuccess("items").BindFunc(func(e *core.RecordEvent) error {
        // send notification, update cache, etc.
        return e.Next()
    })
}
```

## Cron Jobs

```go
func (s *AppServer) registerCronJobs() {
    s.App.Cron().MustAdd("cleanup old records", "0 * * * *", func() {
        // runs every hour
    })

    s.App.Cron().MustAdd("aggregate stats", "*/10 * * * *", func() {
        // runs every 10 minutes
    })
}
```

## Build Tags

Separate dev and production server implementations using gopb SPA helpers:

### Production (`server_production.go`)

```go
//go:build !development

package server

import (
    gopb "github.com/castle-x/go-pocketbase"
    "github.com/pocketbase/pocketbase/core"
    "your-project/site"
)

func (s *AppServer) serveFrontend(se *core.ServeEvent) {
    gopb.ServeSPA(se, site.DistDirFS, []string{"/assets/", "/static/"})
}
```

### Development (`server_development.go`)

```go
//go:build development

package server

import (
    gopb "github.com/castle-x/go-pocketbase"
    "github.com/pocketbase/pocketbase/core"
)

func (s *AppServer) serveFrontend(se *core.ServeEvent) {
    gopb.ServeDevProxy(se, "localhost:5173")
}
```

Run dev mode: `go run -tags development ./cmd/app serve`

**Note:** `gopb.ServeSPA` and `gopb.ServeDevProxy` are standalone functions (not methods). They register a catch-all `GET /{path...}` route, so call them AFTER all other route registrations.

## Custom CLI Commands

PocketBase uses Cobra. Add commands to `app.RootCmd`:

```go
app.RootCmd.AddCommand(&cobra.Command{
    Use:   "seed",
    Short: "Seed database with sample data",
    Run: func(cmd *cobra.Command, args []string) {
        // seeding logic
    },
})
```

The default command is `serve`. Run with: `./app serve --http=0.0.0.0:8090`

## Settings Configuration

```go
func (s *AppServer) initialize(e *core.ServeEvent) error {
    settings := e.App.Settings()
    settings.Meta.AppName = "My App"
    settings.Meta.HideControls = true
    settings.Batch.Enabled = true
    settings.Logs.MinLevel = 4  // warn+
    return e.App.Save(settings)
}
```

## Pitfalls & Resolved Issues

### 1. Password validation — use `record.ValidatePassword()`

PocketBase stores passwords as bcrypt hashes internally. Accessing the hash via `record.GetString("password")` returns an **empty string** — the field is write-only.

```go
// WRONG — always fails, GetString returns ""
hash := record.GetString("password")
bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))

// CORRECT — PocketBase handles bcrypt internally
record.ValidatePassword("plaintext") // returns bool
```

`gopb.IsDefaultPassword(record)` uses this correctly.

### 2. Dual collection initialization (`_superusers` + `users`)

PocketBase has two separate auth collections:

| Collection | Purpose | Login method |
|------------|---------|--------------|
| `_superusers` | Admin UI at `/_/` | Admin login page only |
| `users` | Frontend app login | `pb.collection("users").authWithPassword()` |

**`_superusers` does NOT support public API authentication.** Frontend code cannot call `pb.collection("_superusers").authWithPassword()` — it will fail silently or return an error.

Always create **both** accounts with the same credentials:
```go
s.EnsureDefaults() // creates both if their collections are empty
```

### 3. Password sync between collections

When a user changes their password through the app, only the `users` record is updated by default. The `_superusers` record keeps the old password, causing confusion when accessing the Admin UI.

Solution (built into `gopb`): the `handleChangePassword` endpoint syncs the new password to any `_superusers` records that still use the default password.

### 4. Login state isolation

The `_superusers` and `users` collections have **independent auth sessions**. Logging into the app does NOT log you into the Admin UI. This is a PocketBase architectural characteristic, not a bug.

Current mitigation: password consistency (same credentials work for both). Session bridging would require custom token exchange — not implemented.

### 5. `e.Auth` in protected routes

In routes protected by `apis.RequireAuth()`, `e.Auth` is the authenticated `*core.Record`:

```go
e.Auth.Id           // user ID
e.Auth.Email()      // email
e.Auth.GetString("role")  // custom field
e.Auth.ValidatePassword("test") // password check
```

`e.Auth` is `nil` if the route is not behind `apis.RequireAuth()` — always check.
