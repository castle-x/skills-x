# Go Module Dependencies

## Recommended go.mod

This architecture uses **only Go standard library** - no external HTTP framework dependencies required!

```go
module your-project

go 1.21
```

That's it! No external dependencies needed for the core functionality.

## Go Version Requirements

- **Minimum:** Go 1.16 (embed package introduced)
- **Recommended:** Go 1.21+ (best performance and security)

## Standard Library Dependencies

All imports are from Go standard library:

```go
import (
    "embed"      // Go 1.16+ - for embedding files
    "io/fs"      // Go 1.16+ - filesystem interface
    "net/http"   // Standard HTTP server/client
    "path"       // Path manipulation
    "strings"    // String utilities
)
```

## Version Compatibility

| Go Version | embed Support | Recommended |
|------------|---------------|-------------|
| Go 1.16+ | ✅ Basic support | Minimum |
| Go 1.18+ | ✅ Generics | Good |
| Go 1.21+ | ✅ Best performance | **Recommended** |
| Go 1.22+ | ✅ Enhanced routing | Latest |

## Quick Setup Commands

```bash
# Initialize module
go mod init your-project

# No external dependencies needed!
# Just ensure Go 1.16+

# Verify Go version
go version

# Tidy (will be minimal since no external deps)
go mod tidy
```

## Why Standard Library Only?

| Benefit | Description |
|---------|-------------|
| 🎯 Zero Dependencies | No version conflicts, no supply chain risks |
| 📦 Smaller Binary | No framework overhead |
| 🔒 Security | Only battle-tested stdlib code |
| 🚀 Stability | No breaking changes from external libs |
| 📖 Simplicity | Easy to understand and maintain |

## Optional: If You Need a Framework

If your project requires additional features (middleware, routing groups, etc.), you can optionally add a framework:

```go
// Optional - Gin (most popular)
require github.com/gin-gonic/gin v1.10.0

// Optional - Chi (lightweight, stdlib compatible)
require github.com/go-chi/chi/v5 v5.0.12

// Optional - Echo (minimalist)
require github.com/labstack/echo/v4 v4.12.0
```

But for basic SPA embedding, **Go standard library is sufficient and recommended**.
