# Frontend Patterns

## Table of Contents

1. [PocketBase JS SDK Setup](#pocketbase-js-sdk-setup)
2. [Authentication](#authentication)
3. [Auth State Management](#auth-state-management)
4. [SPA Embedding (embed.go)](#spa-embedding)
5. [Production Static Serving](#production-static-serving)
6. [Development Reverse Proxy](#development-reverse-proxy)
7. [Vite Configuration](#vite-configuration)
8. [PocketBase API Usage](#pocketbase-api-usage)

## PocketBase JS SDK Setup

Install:

```bash
pnpm add pocketbase
```

Create client — connects to **same origin** (no separate port):

```typescript
import PocketBase from "pocketbase";

export const pb = new PocketBase("/");
// In production, frontend is served by the same Go binary
// No proxy needed — all PocketBase APIs available at same origin
```

If using subpath deployment (e.g., `/app/`), pass the base path:

```typescript
const basePath = window.__BASE_PATH__ || "";
export const pb = new PocketBase(basePath);
```

## Authentication

### Login with email/password

```typescript
try {
  const authData = await pb.collection("users").authWithPassword(email, password);
  // authData.token — JWT token
  // authData.record — user record { id, email, ... }
} catch (err) {
  // handle error
}
```

### OAuth2 login

```typescript
await pb.collection("users").authWithOAuth2({
  provider: "github",
});
```

### Auth refresh (verify token is still valid)

```typescript
export const verifyAuth = () => {
  pb.collection("users")
    .authRefresh()
    .catch(() => {
      pb.authStore.clear();
      // redirect to login
    });
};
```

### Logout

```typescript
export function logOut() {
  pb.authStore.clear();
  pb.realtime.unsubscribe();
}
```

### Check auth state

```typescript
pb.authStore.isValid    // boolean — token exists and not expired
pb.authStore.token      // string — JWT token
pb.authStore.record     // RecordModel — current user record
pb.authStore.record?.id // string — user ID

// Listen for auth changes
pb.authStore.onChange((token, record) => {
  if (!token) {
    // logged out
  }
});
```

### Helper functions

```typescript
export const isAdmin = () => pb.authStore.record?.role === "admin";
export const isReadOnly = () => pb.authStore.record?.role === "readonly";
```

## Auth State Management

### With Zustand

```typescript
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { pb } from "./api";

interface AuthState {
  isAuthenticated: boolean;
  setAuthenticated: (v: boolean) => void;
}

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      isAuthenticated: pb.authStore.isValid,
      setAuthenticated: (v) => set({ isAuthenticated: v }),
    }),
    { name: "auth" }
  )
);
```

### With nanostores (lighter alternative)

```typescript
import { atom } from "nanostores";
export const $authenticated = atom(pb.authStore.isValid);
```

### Sync with PocketBase authStore

```typescript
pb.authStore.onChange(() => {
  if (!pb.authStore.isValid) {
    useAuth.getState().setAuthenticated(false);
    // or: $authenticated.set(false)
  }
});
```

## SPA Embedding

`site/embed.go`:

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

## Production Static Serving

Use `gopb.ServeSPA` from `github.com/castle-x/go-pocketbase`:

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

`gopb.ServeSPA` handles:
- Reading `index.html` for SPA fallback
- Serving static assets with `Cache-Control: public, max-age=2592000`
- Catch-all route returning `index.html` for client-side routing

## Development Reverse Proxy

Use `gopb.ServeDevProxy` to proxy to Vite dev server:

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

Run with: `go run -tags development ./cmd/app serve`

This gives you:
- PocketBase API + custom routes at `:8090`
- Vite HMR via the reverse proxy (also at `:8090`)
- Single origin — no CORS issues

## Vite Configuration

In dev mode, Vite can also proxy API calls to Go backend (alternative to Go reverse proxy above):

```typescript
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8090",
        changeOrigin: true,
      },
      "/_": {
        target: "http://localhost:8090",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
```

When using the Go reverse proxy approach (`-tags development`), the Vite proxy is not needed — Go proxies Vite instead.

## Setup Status Integration

The backend (gopb) provides setup endpoints for first-run password change. Frontend should check after login:

```typescript
// Check if user needs to change default password
const { needsPasswordChange } = await api.get("setup/status").json();
if (needsPasswordChange) {
  // Show forced password change dialog
}

// Change password (syncs to both users and _superusers)
await api.post("setup/change-password", {
  json: { password: newPwd, passwordConfirm: newPwd }
});

// Re-authenticate with new password to refresh token
await pb.collection("users").authWithPassword(email, newPwd);
```

**Key points:**
- Login ONLY against `users` collection — `_superusers` does not support public API auth
- After password change, re-authenticate to get a fresh token
- The `Authorization` header for custom API routes uses `pb.authStore.token`

## PocketBase API Usage

### CRUD via JS SDK (auto-handles auth token)

```typescript
// List with pagination
const result = await pb.collection("items").getList(1, 20, {
  filter: 'status = "active"',
  sort: "-created",
});

// Get one
const record = await pb.collection("items").getOne("RECORD_ID");

// Create
const record = await pb.collection("items").create({
  title: "Hello",
  content: "World",
});

// Update
await pb.collection("items").update("RECORD_ID", { title: "Updated" });

// Delete
await pb.collection("items").delete("RECORD_ID");
```

### Custom API endpoints

```typescript
// GET
const data = await pb.send("/api/v1/custom-endpoint", {});

// POST with body
const data = await pb.send("/api/v1/custom-endpoint", {
  method: "POST",
  body: JSON.stringify({ key: "value" }),
});
```

### Realtime subscriptions

```typescript
pb.collection("items").subscribe("*", (e) => {
  console.log(e.action); // "create" | "update" | "delete"
  console.log(e.record);
});

// Unsubscribe
pb.collection("items").unsubscribe("*");

// Unsubscribe all
pb.realtime.unsubscribe();
```

### File uploads

```typescript
const formData = new FormData();
formData.append("title", "My File");
formData.append("document", fileInput.files[0]);

const record = await pb.collection("items").create(formData);

// Get file URL
const url = pb.files.getURL(record, record.document);
```
