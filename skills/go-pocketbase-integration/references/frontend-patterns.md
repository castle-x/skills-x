# Frontend Patterns

## Table of Contents

1. [pb-client.ts — Auth Center](#pb-clientts--auth-center)
2. [Authentication](#authentication)
3. [Auth State Management](#auth-state-management)
4. [With External API Clients](#with-external-api-clients)
5. [SPA Embedding (embed.go)](#spa-embedding)
6. [Production Static Serving](#production-static-serving)
7. [Development Reverse Proxy](#development-reverse-proxy)
8. [Vite Configuration](#vite-configuration)
9. [PocketBase API Usage](#pocketbase-api-usage)

## pb-client.ts — Auth Center

`pb-client.ts` is more than a PocketBase instance creator — it is the **auth center** for the entire frontend. It is the dependency root for all auth-related modules, loaded before any React component or Zustand store.

Install:

```bash
pnpm add pocketbase
```

Full `pb-client.ts` pattern (`src/shared/lib/pb-client.ts`):

```typescript
import PocketBase from "pocketbase";

// ═══════════════════════════════════════════════════
// 1. PocketBase instance
// ═══════════════════════════════════════════════════

export const pb = new PocketBase("/");
// In production, frontend is served by the same Go binary.
// No proxy needed — all PocketBase APIs available at same origin.
//
// If using subpath deployment (e.g., /app/):
//   const basePath = window.__BASE_PATH__ || "";
//   export const pb = new PocketBase(basePath);

// ═══════════════════════════════════════════════════
// 2. Auth state restoration
// ═══════════════════════════════════════════════════
//
// Restore auth state from localStorage immediately on module load.
// This ensures pb.authStore.token is ready before any other module reads it.
//
// Why here (not in React components or Zustand store)?
// - pb-client.ts is the earliest-loaded module (all pb users import it)
// - Auth restoration here guarantees the earliest possible timing
//
// AUTH_STORAGE_KEY must match the Zustand persist `name`.
// If not using Zustand, use your project's own storage key and adjust
// the JSON parsing logic accordingly.

const AUTH_STORAGE_KEY = "auth-storage";

function restoreAuthFromStorage(): void {
  if (typeof window === "undefined") return;

  const stored = localStorage.getItem(AUTH_STORAGE_KEY);
  if (!stored) return;

  try {
    const { state } = JSON.parse(stored) as {
      state?: { token?: string; user?: unknown };
    };
    if (state?.token && state?.user) {
      pb.authStore.save(state.token, state.user);
    }
  } catch (e) {
    console.error("[pb-client] Failed to restore auth state:", e);
  }
}

restoreAuthFromStorage();

// ═══════════════════════════════════════════════════
// 3. Auth-aware fetch factory
// ═══════════════════════════════════════════════════
//
// Returns a fetch function that lazily reads pb.authStore.token on each call
// and injects it as the Authorization header.
//
// Lazy (not a closure snapshot): token refresh is automatically reflected
// without rebuilding the fetch function.
//
// Return type is `typeof globalThis.fetch` — compatible with any library
// that accepts a custom fetch (ofetch, ky, generated RPC clients, etc.).

export function createAuthFetch(): typeof globalThis.fetch {
  return (input, init) => {
    const headers = new Headers(init?.headers);
    const token = pb.authStore.token;
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
    return globalThis.fetch(input, { ...init, headers });
  };
}

// ═══════════════════════════════════════════════════
// 4. Error handler factory
// ═══════════════════════════════════════════════════
//
// Returns a unified error handler. The primary case is 401 Unauthorized:
// clears auth state, which triggers pb.authStore.onChange → UI logout.
//
// Error message format convention: "request failed: status={code}, body={text}"
// If your HTTP client uses a different error format, create a project-side
// wrapper rather than modifying pb-client.ts.

export function createErrorHandler(): (error: Error, method: string) => void {
  return (error, method) => {
    if (error.message.includes("status=401")) {
      console.warn(`[pb-client] 401 on ${method}, clearing auth`);
      pb.authStore.clear();
    }
  };
}

// ═══════════════════════════════════════════════════
// 5. Convenience helpers
// ═══════════════════════════════════════════════════

export function isAuthenticated(): boolean {
  return pb.authStore.isValid;
}

export function getCurrentUser() {
  return pb.authStore.record;
}
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

`pb-client.ts` handles auth restoration. The Zustand store manages UI state only, keeping its responsibility narrow.

### useAuth.ts (Zustand — UI state only)

Auth restoration is done by `pb-client.ts`. `useAuth.ts` does not need restoration logic.

```typescript
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { pb } from "@/shared/lib/pb-client";

interface AuthUser {
  id: string;
  email: string;
  [key: string]: unknown;
}

interface AuthState {
  user: AuthUser | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: AuthUser, token: string) => void;
  logout: () => void;
}

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      setAuth: (user, token) => {
        pb.authStore.save(token, user);   // sync to PocketBase
        set({ user, token, isAuthenticated: true });
      },

      logout: () => {
        pb.authStore.clear();             // sync to PocketBase
        set({ user: null, token: null, isAuthenticated: false });
      },
    }),
    { name: "auth-storage" },            // ← must match AUTH_STORAGE_KEY in pb-client.ts
  ),
);
```

### providers.tsx — Auth guard (React)

Sync `pb.authStore` events to Zustand. `pb.authStore.clear()` (triggered by 401 handler or logout) fires `onChange`, which calls `zustand.logout()`. No recursion risk: once `isAuthenticated` is false, the `if` condition is not met.

```typescript
import { useEffect } from "react";
import { pb } from "@/shared/lib/pb-client";
import { useAuth } from "@/shared/hooks/useAuth";

function AuthGuard({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    const unsubscribe = pb.authStore.onChange(() => {
      const zustand = useAuth.getState();
      if (!pb.authStore.isValid && zustand.isAuthenticated) {
        zustand.logout();
        // navigate to login page here
      }
    });
    return () => unsubscribe();
  }, []);

  return <>{children}</>;
}
```

### With nanostores (lighter alternative)

```typescript
import { atom } from "nanostores";
import { pb } from "@/shared/lib/pb-client";

export const $authenticated = atom(pb.authStore.isValid);

pb.authStore.onChange(() => {
  $authenticated.set(pb.authStore.isValid);
});
```

### Auth lifecycle summary

| Event | Who acts | What happens |
|-------|----------|--------------|
| Page load | `pb-client.ts` | `restoreAuthFromStorage()` → `pb.authStore.save(token, user)` |
| Login | App code | `pb.collection("users").authWithPassword()` → `useAuth.setAuth()` → persist |
| Logout | `useAuth.logout()` | `pb.authStore.clear()` → Zustand cleared → localStorage cleared |
| 401 response | `createErrorHandler()` | `pb.authStore.clear()` → `onChange` → `AuthGuard` → `useAuth.logout()` |

## With External API Clients

`createAuthFetch()` returns `typeof globalThis.fetch` — inject it into any library that accepts a custom fetch.

### Direct use

```typescript
import { createAuthFetch } from "@/shared/lib/pb-client";

const authFetch = createAuthFetch();

const resp = await authFetch("/api/some-endpoint", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ key: "value" }),
});
```

### With ofetch or ky

```typescript
// ofetch
import { ofetch } from "ofetch";
import { createAuthFetch } from "@/shared/lib/pb-client";

const api = ofetch.create({ fetch: createAuthFetch() });

// ky
import ky from "ky";
import { createAuthFetch } from "@/shared/lib/pb-client";

const api = ky.create({ fetch: createAuthFetch() });
```

### With generated RPC clients

If the project uses a code-generation tool that produces typed clients with injectable fetch:

```typescript
import { createAuthFetch, createErrorHandler } from "@/shared/lib/pb-client";
import { SomeGeneratedClient } from "@/api/generated/client";

// pb-client.ts does not depend on the generated client —
// it only provides standard-typed building blocks.
const client = new SomeGeneratedClient("/api/v1", {
  fetch: createAuthFetch(),
  onError: createErrorHandler(),
});
```

> `createErrorHandler()` detects 401 via `error.message.includes("status=401")`. If your HTTP client formats errors differently, wrap it on the project side rather than modifying `pb-client.ts`.

### Composing multiple fetch wrappers

```typescript
import { createAuthFetch } from "@/shared/lib/pb-client";

const authFetch = createAuthFetch();

const tracedFetch: typeof fetch = (input, init) => {
  const headers = new Headers(init?.headers);
  headers.set("X-Request-ID", crypto.randomUUID());
  return authFetch(input, { ...init, headers });
};
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

Use `gopb.ServeSPA` from `github.com/castle-x/goutils/pocketbase`:

```go
//go:build !development

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
    gopb "github.com/castle-x/goutils/pocketbase"
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
