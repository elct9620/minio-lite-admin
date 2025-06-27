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

- **HTTP Router**: Chi router (`github.com/go-chi/chi/v5`) with custom zerolog middleware (Recoverer, RequestID)
- **Configuration**: Viper-based config management with flags, environment variables, and YAML support
- **Logging**: Zerolog structured logging with configurable levels and pretty printing
- **Handler Structure**: Organized in `internal/handler/http` package with separation of concerns
- **API Structure**: RESTful endpoints under `/api` prefix with structured JSON responses
- **Integration**: Uses `github.com/olivere/vite` for HTML fragment generation and asset serving
- **Deployment**: Single binary with embedded frontend assets

### Project Structure

```
internal/
├── config/          # Configuration management with Viper
│   └── config.go    # Config structs and loading logic
├── infra/           # Infrastructure layer (factories)
│   └── minio.go     # MinIO admin client factory
├── logger/          # Zerolog configuration and setup
│   └── logger.go    # Logger initialization and level parsing
├── service/         # Business logic layer
│   └── get_server_info_service.go # MinIO server info service
└── handler/
    └── http/        # HTTP handlers
        ├── api.go   # API endpoints (health, server-info, minio/server-info)
        ├── frontend.go # Frontend serving with Vite integration
        └── middleware.go # Custom chi middleware (zerolog request logger)
```

### Configuration Management

Uses Viper for robust configuration with multiple sources (priority order):
1. **Command line flags** (highest priority)
2. **Environment variables** (`MINIO_ADMIN_*` prefix)
3. **Configuration files** (YAML)
4. **Default values** (lowest priority)

**Environment Variables**:
- `MINIO_ADMIN_ADDR`: Server address (default: `:8080`)
- `MINIO_ADMIN_DEV`: Development mode (default: `false`)
- `MINIO_ADMIN_VITE_URL`: Vite dev server URL (default: `http://localhost:5173`)
- `MINIO_ADMIN_VITE_ENTRY`: Vite entry point (default: `/src/main.ts`)
- `MINIO_ADMIN_LOG_LEVEL`: Log level (default: `info`, options: trace, debug, info, warn, error, fatal, panic)
- `MINIO_ADMIN_LOG_PRETTY`: Pretty log output (default: `true`)
- `MINIO_URL`: MinIO server URL (default: `http://localhost:9000`)
- `MINIO_ROOT_USER`: MinIO root username (required)
- `MINIO_ROOT_PASSWORD`: MinIO root password (required)

### Handler Architecture

- **API Handlers**: Type-safe JSON responses with proper error handling
- **Frontend Handler**: Dependency injection pattern with config and embedded assets
- **Custom Middleware**: Zerolog-based request logging with structured output (method, URI, status, elapsed time)
- **Separation of Concerns**: Clean boundaries between routing, configuration, and business logic

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
# Development mode (preferred for development)
go run ./main.go -dev

# Production mode (serves embedded assets)
go run ./main.go

# Custom port
go run ./main.go -addr :3000

# Build binary for verification (outputs to bin/)
go build . -o bin/minio-lite-admin
./bin/minio-lite-admin -dev
```

### Full Development Workflow
1. Terminal 1: `pnpm dev` (starts Vite dev server)
2. Terminal 2: `go run ./main.go -dev` (starts Go server in dev mode)
3. Visit http://localhost:8080

### Docker Development
```bash
# Start full development environment with watch mode
docker compose up --watch

# Build production Docker image
docker build -t minio-lite-admin .
```

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

**Important**: The `ViteURL` must be accessible by the **browser**, not the Go server. The `olivere/vite` package generates HTML with script tags pointing to this URL - it does NOT proxy Vite assets. 

**Asset Serving Strategy**:
- **Development**: All static assets (JS, CSS, images) are served by Vite dev server at localhost:5173
- **Production**: All assets are embedded in Go binary after `vite build` creates the `dist/` directory

In Docker development:
- Use `VITE_URL=http://localhost:5173` (browser-accessible)
- NOT `http://frontend:5173` (only accessible inside Docker network)

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
- ✅ Hybrid Go-Vue integration with `olivere/vite`
- ✅ Chi router with custom zerolog middleware (Recoverer, RequestID)
- ✅ Viper configuration management (flags, env vars, config files)
- ✅ Zerolog structured logging with configurable levels and pretty printing
- ✅ MinIO Admin SDK integration with service layer architecture
- ✅ Infrastructure layer with MinIO client factory (`internal/infra`)
- ✅ Service layer for business logic (`internal/service`)
- ✅ Context-based logging for request tracing
- ✅ `/api/minio/server-info` endpoint with MinIO integration
- ✅ Structured HTTP handlers in `internal/handler/http`
- ✅ Vue.js 3 + TypeScript frontend scaffold
- ✅ Docker development environment with watch mode and MinIO service
- ✅ Production Docker build with multi-stage process
- ✅ Development/production mode switching

**Next Steps (TODO)**:
- Additional MinIO administrative features (disk status, access keys, replication)
- Frontend components for MinIO management UI
- MinIO connection health check and validation
- Authentication and authorization system
- More API endpoints for MinIO operations (buckets, users, policies)

## License Constraints

This project uses AGPLv3 license due to dependency on `github.com/minio/madmin-go` which is AGPLv3 licensed. Any derivative work must maintain AGPLv3 compatibility.

## Development Notes

- The project requires Go 1.24+ and Node.js for development
- Frontend assets must be built (`pnpm build`) before production deployment
- The `dist/.keep` file ensures the dist directory exists for Go embed
- API endpoints are structured under `/api` prefix for clear separation
- Static assets are served by Vite dev server in development, embedded in Go binary in production
- Build output goes to `bin/` directory (excluded from git) for verification without cleanup
- Use `go run ./main.go -dev` for development, `go build . -o bin/minio-lite-admin` for verification
- Use `gofmt -w .` to format all Go files before running verification tests