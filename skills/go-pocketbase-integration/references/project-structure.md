# Project Structure

## Directory Layout

```
project/
├── go.mod                          # Go module (requires pocketbase dependency)
├── go.sum
├── Makefile                        # Build, dev, deploy commands
├── cmd/
│   └── app/
│       └── main.go                 # Entry point: create PocketBase + start AppServer
├── internal/
│   ├── server/
│   │   ├── server.go               # AppServer struct wrapping core.App, OnServe hook
│   │   ├── server_production.go    # //go:build !development — serve embedded SPA
│   │   └── server_development.go   # //go:build development — reverse proxy to Vite
│   ├── migrations/
│   │   ├── initial_settings.go     # First-run setup (settings, superuser)
│   │   └── 0_collections.go        # Collection schema snapshot
│   └── {feature}/                  # Feature packages (users, records, etc.)
│       └── {feature}.go
├── site/                           # Frontend SPA
│   ├── embed.go                    # go:embed directive
│   ├── src/
│   │   ├── lib/
│   │   │   └── api.ts              # PocketBase JS client: new PocketBase(basePath)
│   │   ├── components/
│   │   │   └── login/              # Auth forms using pb.collection("users").authWithPassword()
│   │   └── ...
│   ├── dist/                       # Build output (embedded into binary)
│   ├── package.json                # pocketbase JS SDK dependency
│   └── vite.config.ts
└── Dockerfile                      # Multi-stage: build → scratch
```

## Key Files

### `cmd/app/main.go`

- Creates `pocketbase.NewWithConfig()` with `DefaultDataDir` and `DefaultDev`
- Registers migration command via `migratecmd.MustRegister()`
- Adds custom CLI commands via `app.RootCmd.AddCommand()`
- Passes app to AppServer, calls `srv.Start()`

### `internal/server/server.go`

- `AppServer` struct embeds `*gopb.AppServer` — gains all PocketBase + gopb methods
- `Start()` registers OnServe hook, then calls `s.AppServer.Start()`
- Registration order in OnServe: setup routes → business routes → ensure defaults → serve frontend

### `site/embed.go`

```go
package site

import (
    "embed"
    "io/fs"
)

//go:embed all:dist
var distDir embed.FS

var DistDirFS, _ = fs.Sub(distDir, "dist")
```

### `internal/migrations/`

- Blank-imported in main: `_ "your-project/internal/migrations"`
- Each file has `func init()` calling `m.Register(upFunc, downFunc)`
- Auto-generated when `Automigrate: true` and collections change in admin UI

## go.mod Dependencies

Core required:

```
github.com/castle-x/goutils/pocketbase latest  # gopb — AppServer, defaults, setup routes, SPA helpers
github.com/pocketbase/pocketbase  v0.36+  # PocketBase core
github.com/pocketbase/dbx                 # included transitively, used for raw queries
```

PocketBase brings in: SQLite (modernc.org/sqlite), JWT, imaging, mail, validation, etc.

## Makefile Targets

```makefile
PB_DATA := app_data

dev:          # go run with -tags development (proxies to Vite)
dev-web:      # cd site && pnpm dev
build-web:    # cd site && pnpm build
build:        # build-web + go build (embeds dist/)
build-all:    # cross-platform builds
clean:        # rm -rf site/dist bin/
```

Build order: **frontend first, then backend** (so go:embed picks up dist/).
