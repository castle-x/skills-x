# Embed Implementation

This reference provides the complete `embed.go` implementation for embedding frontend static files into Go binary.

## File: `site/embed.go`

```go
package site

import (
	"embed"
	"io/fs"
)

// distDir embeds all files from the dist directory
// The `all:` prefix includes hidden files (starting with . or _)
//
//go:embed all:dist
var distDir embed.FS

// DistDirFS is the public filesystem interface for embedded files
// fs.Sub removes the "dist/" prefix, so files are accessed directly:
//   - "index.html" instead of "dist/index.html"
//   - "assets/index-abc123.js" instead of "dist/assets/index-abc123.js"
var DistDirFS, _ = fs.Sub(distDir, "dist")
```

---

## Directory Structure

The `embed.go` file must be in the same directory as the `dist/` folder:

```
site/
├── embed.go          ← This file
├── dist/             ← Build output (to be embedded)
│   ├── index.html
│   ├── assets/
│   │   ├── index-abc123.js
│   │   └── index-xyz789.css
│   └── favicon.ico
├── src/              ← Source files (NOT embedded)
├── package.json
└── vite.config.ts
```

---

## Embed Directive Syntax

### Basic Syntax

```go
//go:embed <pattern>
var <variable> embed.FS
```

### Pattern Options

| Pattern | Description |
|---------|-------------|
| `dist` | Embeds `dist/` directory (excludes hidden files) |
| `all:dist` | Embeds `dist/` directory **including hidden files** |
| `dist/*` | Embeds only direct children of `dist/` |
| `dist/**` | Invalid - use `all:dist` instead |

### Why Use `all:dist`?

The `all:` prefix ensures hidden files are included. Some build tools may generate:
- `.htaccess` - Apache config
- `_redirects` - Netlify redirects
- `.well-known/` - ACME challenges

Without `all:`, these would be silently excluded.

---

## fs.Sub() Explanation

```go
var DistDirFS, _ = fs.Sub(distDir, "dist")
```

This creates a sub-filesystem that removes the `dist/` prefix:

| Without fs.Sub | With fs.Sub |
|----------------|-------------|
| `distDir.Open("dist/index.html")` | `DistDirFS.Open("index.html")` |
| `distDir.Open("dist/assets/main.js")` | `DistDirFS.Open("assets/main.js")` |

This makes the embedded filesystem work identically to a real directory structure.

---

## Usage in Other Packages

### Import and Use

```go
package main

import (
	"net/http"
	
	"your-project/site"
)

func main() {
	// site.DistDirFS is available as fs.FS interface
	fileServer := http.FileServer(http.FS(site.DistDirFS))
	http.Handle("/", fileServer)
	http.ListenAndServe(":8080", nil)
}
```

### Read Specific File

```go
import (
	"io/fs"
	
	"your-project/site"
)

func readIndex() ([]byte, error) {
	return fs.ReadFile(site.DistDirFS, "index.html")
}
```

### List Files

```go
import (
	"io/fs"
	
	"your-project/site"
)

func listFiles() {
	fs.WalkDir(site.DistDirFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
}
```

---

## Build Requirements

### Build Order

**Critical:** Frontend must be built BEFORE Go compilation!

```bash
# 1. Build frontend first
cd site && npm run build   # Creates site/dist/

# 2. Then build Go binary
go build ./cmd/app         # Embeds site/dist/ into binary
```

If `dist/` doesn't exist during `go build`, compilation will fail.

### Empty dist/ Handling

For development without frontend, create a minimal placeholder:

```bash
mkdir -p site/dist
echo '<!DOCTYPE html><html><body>Frontend not built</body></html>' > site/dist/index.html
```

---

## Alternative Patterns

### Multiple Directories

```go
//go:embed all:dist all:static
var content embed.FS
```

### Specific Files

```go
//go:embed dist/index.html dist/favicon.ico
var content embed.FS
```

### Pattern Matching

```go
//go:embed dist/*.html dist/assets/*.js
var content embed.FS
```

---

## Common Errors

### Error: pattern dist: no matching files found

**Cause:** `dist/` directory doesn't exist

**Solution:**
```bash
cd site && npm run build
```

### Error: embed: invalid pattern syntax

**Cause:** Invalid glob pattern like `dist/**`

**Solution:** Use `all:dist` instead

### Runtime Error: file does not exist

**Cause:** Accessing wrong path (missing or extra `dist/` prefix)

**Solution:** Use `fs.Sub()` to normalize paths
