# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MinIO Lite Admin is a lightweight web-based administration tool for MinIO instances, built with a hybrid Go backend + Vue.js frontend architecture. The project addresses the gap created by MinIO's removal of management features from the community edition, focusing on essential administrative functions:

- Disk Status Monitoring
- Access Key Management  
- Site Replication Configuration

## Architecture

### Hybrid Go-Vue Integration Pattern

The project uses `github.com/olivere/vite` for seamless integration between Go backend and Vue.js frontend:

- **Development Mode**: Go server (`main.go -dev`) proxies to Vite dev server on localhost:5173
- **Production Mode**: Go server serves embedded assets from `dist/` directory using Go's `embed` directive
- **Template System**: HTML templates with Vite fragment injection (`{{ .Vite.Tags }}`) for proper asset loading

### Backend Architecture

- **HTTP Router**: Chi router (`github.com/go-chi/chi/v5`) with middleware (Logger, Recoverer, RequestID)
- **API Structure**: RESTful endpoints under `/api` prefix
- **Integration**: Uses `github.com/olivere/vite` for HTML fragment generation and asset serving
- **Deployment**: Single binary with embedded frontend assets

### Frontend Architecture

- **Framework**: Vue.js 3 with Composition API and `<script setup>` SFCs
- **Build Tool**: Vite with TypeScript support
- **Package Manager**: pnpm (evidenced by `pnpm-lock.yaml`)
- **Configuration**: Vite configured to generate manifest.json for Go integration

## Development Commands

### Frontend Development
```bash
# Install dependencies
pnpm install

# Start Vite dev server (port 5173)
pnpm dev

# Type check and build for production
pnpm build

# Preview production build
pnpm preview
```

### Backend Development
```bash
# Development mode (requires Vite dev server running)
go run main.go -dev

# Production mode (serves embedded assets)
go run main.go

# Custom port
go run main.go -addr :3000
```

### Full Development Workflow
1. Terminal 1: `pnpm dev` (starts Vite dev server)
2. Terminal 2: `go run main.go -dev` (starts Go server in dev mode)
3. Visit http://localhost:8080

## Key Implementation Details

### Vite Integration Pattern

The project uses the HTML fragment approach rather than the HTTP handler approach:

```go
// Template with Vite fragment injection
const indexTemplate = `<!doctype html>
<html lang="en">
  <head>
    <title>MinIO Lite Admin</title>
    {{ .Vite.Tags }}
  </head>
  <body>
    <div id="app"></div>
  </body>
</html>`

// Fragment generation for dev/prod modes
viteFragment, err := vite.HTMLFragment(vite.Config{
    FS:        os.DirFS("."),
    IsDev:     *isDev,
    ViteURL:   "http://localhost:5173",
    ViteEntry: "/src/main.ts",
})
```

### Build Configuration

**Vite Configuration** (`vite.config.ts`):
```typescript
export default defineConfig({
  plugins: [vue()],
  build: {
    manifest: true,        // Required for Go integration
    rollupOptions: {
      input: '/src/main.ts' // Entry point for build
    }
  }
})
```

### Current State and Next Steps

**Completed**:
- Basic Go-Vite integration setup
- Chi router with middleware
- Vue.js application scaffold
- Development/production mode switching

**Missing (TODO)**:
- MinIO Admin SDK integration (`github.com/minio/madmin-go`)
- Actual MinIO administrative features
- Frontend components for MinIO management
- Configuration system for MinIO connections
- Authentication and authorization

## License Constraints

This project uses AGPLv3 license due to dependency on `github.com/minio/madmin-go` which is AGPLv3 licensed. Any derivative work must maintain AGPLv3 compatibility.

## Development Notes

- The project requires Go 1.24+ and Node.js for development
- Frontend assets must be built (`pnpm build`) before production deployment
- The `dist/.keep` file ensures the dist directory exists for Go embed
- API endpoints are structured under `/api` prefix for clear separation
- Static assets are served differently in dev (filesystem) vs prod (embedded) modes