# Vite Configuration

## File Location

`site/vite.config.ts`

## Complete Configuration

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  
  build: {
    // Output directory - MUST match go:embed path
    outDir: 'dist',
    
    // Clean dist before build
    emptyDirBeforeWrite: true,
    
    // Generate source maps for debugging (optional)
    sourcemap: false,
    
    // Chunk size warning limit
    chunkSizeWarningLimit: 1000,
  },
  
  server: {
    // Dev server port
    port: 5173,
    
    // Proxy configuration for development
    proxy: {
      // Proxy API requests to Go backend
      '/apis': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      
      // Proxy WebSocket connections
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
})
```

## Vue Version

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  
  build: {
    outDir: 'dist',
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

## Configuration Explained

### Build Options

| Option | Value | Purpose |
|--------|-------|---------|
| `outDir` | `'dist'` | Must match `//go:embed all:dist` path |
| `emptyDirBeforeWrite` | `true` | Clean old files before new build |
| `sourcemap` | `false` | Disable source maps for production |

### Proxy Configuration

The proxy routes requests during development:

```
Browser (localhost:5173)
    │
    ├── /apis/* ──────► Go backend (localhost:8080)
    ├── /ws ──────────► Go backend WebSocket
    └── /* ───────────► Vite dev server (HMR)
```

### Multiple API Paths

```typescript
proxy: {
  '/apis': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
  '/auth': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
  '/upload': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
}
```

### Path Rewriting

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '/apis/v1'),
  },
}
```

## package.json Scripts

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  }
}
```

## Development Workflow

1. Start Go backend: `make dev-backend` (port 8080)
2. Start Vite dev: `make dev-web` (port 5173)
3. Open browser: `http://localhost:5173`
4. API calls proxy to Go backend automatically

## Production Build

```bash
cd site
npm run build
# Output: site/dist/
#   ├── index.html
#   └── assets/
#       ├── index-[hash].js
#       └── index-[hash].css
```
