# Static File Server Implementation

This reference provides a complete implementation for serving embedded SPA files using **Go standard library `net/http`** - no external framework needed!

## File: `pkg/siteserver/siteserver.go`

```go
package siteserver

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// SPAHandler wraps an existing http.Handler and serves embedded static files
// For paths not handled by the wrapped handler, it serves static files from embedFS
// For SPA routes (no file extension), it returns index.html
type SPAHandler struct {
	// apiHandler is the wrapped handler (e.g., http.ServeMux with API routes)
	apiHandler http.Handler
	// embedFS is the embedded filesystem containing static files
	embedFS fs.FS
	// fileServer serves static files from embedFS
	fileServer http.Handler
	// indexHTML is pre-loaded for SPA route fallback
	indexHTML []byte
}

// NewSPAHandler creates a new SPA handler
//
// Parameters:
//   - apiHandler: Handler for API routes (can be nil)
//   - embedFS: Embedded filesystem (e.g., site.DistDirFS)
//
// Returns SPAHandler and error if index.html cannot be read
func NewSPAHandler(apiHandler http.Handler, embedFS fs.FS) (*SPAHandler, error) {
	// Pre-load index.html for SPA route fallback
	// This avoids repeated filesystem reads for every SPA route
	indexHTML, err := fs.ReadFile(embedFS, "index.html")
	if err != nil {
		return nil, err
	}

	return &SPAHandler{
		apiHandler: apiHandler,
		embedFS:    embedFS,
		fileServer: http.FileServer(http.FS(embedFS)),
		indexHTML:  indexHTML,
	}, nil
}

// ServeHTTP implements http.Handler interface
func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := strings.TrimPrefix(r.URL.Path, "/")

	// Try API handler first if it exists
	if h.apiHandler != nil {
		// Check if path starts with known API prefixes
		if strings.HasPrefix(reqPath, "apis/") || strings.HasPrefix(reqPath, "api/") {
			h.apiHandler.ServeHTTP(w, r)
			return
		}
	}

	// Check if request is for a static resource (has file extension)
	if strings.Contains(path.Base(reqPath), ".") {
		// Set long-term cache for assets/ directory (files have hash in name)
		if strings.HasPrefix(reqPath, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		// Try to open file from embedded FS
		if _, err := h.embedFS.Open(reqPath); err == nil {
			// File exists, serve it
			h.fileServer.ServeHTTP(w, r)
			return
		}
	}

	// Not a static file or file doesn't exist
	// Return index.html for SPA client-side routing
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write(h.indexHTML)
}

// WrapHandler is a convenience function to wrap an existing handler
//
// Usage:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/apis/v1/health", healthHandler)
//	handler := siteserver.WrapHandler(mux, site.DistDirFS)
//	http.ListenAndServe(":8080", handler)
func WrapHandler(apiHandler http.Handler, embedFS fs.FS) http.Handler {
	handler, err := NewSPAHandler(apiHandler, embedFS)
	if err != nil {
		// If index.html doesn't exist, just return the API handler
		// This allows the server to start even without frontend build
		return apiHandler
	}
	return handler
}

// ServeStatic creates a handler that only serves static files (no API wrapper)
// Use this when you have a separate API server and only need static file serving
func ServeStatic(embedFS fs.FS) (http.Handler, error) {
	return NewSPAHandler(nil, embedFS)
}
```

---

## Usage Examples

### Example 1: Basic Usage with http.ServeMux

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"your-project/pkg/siteserver"
	"your-project/site"
)

func main() {
	mux := http.NewServeMux()

	// Register API routes
	mux.HandleFunc("/apis/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/apis/v1/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Hello from API",
			"data":    []int{1, 2, 3},
		})
	})

	// Wrap with SPA handler
	handler := siteserver.WrapHandler(mux, site.DistDirFS)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
```

### Example 2: Static Files Only (No API)

```go
package main

import (
	"log"
	"net/http"

	"your-project/pkg/siteserver"
	"your-project/site"
)

func main() {
	handler, err := siteserver.ServeStatic(site.DistDirFS)
	if err != nil {
		log.Fatalf("Failed to create static handler: %v", err)
	}

	log.Println("Static server starting on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
```

### Example 3: With Custom API Prefix Detection

If you need custom API prefix detection, create your own handler:

```go
package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"your-project/site"
)

type CustomSPAHandler struct {
	mux        *http.ServeMux
	fileServer http.Handler
	indexHTML  []byte
	apiPrefixes []string
}

func NewCustomHandler(embedFS fs.FS, apiPrefixes []string) (*CustomSPAHandler, error) {
	indexHTML, err := fs.ReadFile(embedFS, "index.html")
	if err != nil {
		return nil, err
	}

	return &CustomSPAHandler{
		mux:         http.NewServeMux(),
		fileServer:  http.FileServer(http.FS(embedFS)),
		indexHTML:   indexHTML,
		apiPrefixes: apiPrefixes,
	}, nil
}

func (h *CustomSPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := strings.TrimPrefix(r.URL.Path, "/")

	// Check custom API prefixes
	for _, prefix := range h.apiPrefixes {
		if strings.HasPrefix(reqPath, prefix) {
			h.mux.ServeHTTP(w, r)
			return
		}
	}

	// Static file or SPA route handling...
	if strings.Contains(path.Base(reqPath), ".") {
		if strings.HasPrefix(reqPath, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		h.fileServer.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(h.indexHTML)
}

func main() {
	handler, _ := NewCustomHandler(site.DistDirFS, []string{"api/", "v1/", "graphql"})
	http.ListenAndServe(":8080", handler)
}
```

---

## Key Implementation Points

### 1. Static Resource Detection

```go
if strings.Contains(path.Base(reqPath), ".") {
    // Has file extension → static resource
}
```

This checks the **base name** (final path component) for a dot, avoiding false positives like `/api/v1.0/users`.

### 2. SPA Route Fallback

For paths without file extensions (frontend routes like `/dashboard/settings`), return `index.html` and let the frontend router handle it.

### 3. Caching Strategy

| Path | Cache-Control | Reason |
|------|---------------|--------|
| `/assets/*` | `max-age=31536000, immutable` | Files have content hash in filename, safe for long cache |
| `/index.html` | `no-cache` | Entry point must always be fresh to load new assets |
| Other static | Default | Normal browser caching |

### 4. Pre-loading index.html

```go
indexHTML, err := fs.ReadFile(embedFS, "index.html")
```

Reading `index.html` once at startup avoids repeated filesystem reads for every SPA route request. Since it's embedded, this is fast, but caching prevents unnecessary allocations.

### 5. Graceful Degradation

The `WrapHandler` function gracefully handles missing `index.html`:

```go
func WrapHandler(apiHandler http.Handler, embedFS fs.FS) http.Handler {
    handler, err := NewSPAHandler(apiHandler, embedFS)
    if err != nil {
        // Fallback to API-only mode
        return apiHandler
    }
    return handler
}
```

This allows the server to start even if frontend build hasn't been done yet.
