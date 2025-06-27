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
- **Handler Structure**: Service-based architecture in `internal/handler/http` with struct methods and dependency injection
- **API Structure**: RESTful endpoints under `/api` prefix with structured JSON responses
- **Integration**: Uses `github.com/olivere/vite` for HTML fragment generation and asset serving
- **Deployment**: Single binary with embedded frontend assets
- **Build Tags**: Conditional compilation using `dist` tag for production builds (development mode is default)

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
    └── http/        # HTTP handlers (Service-based architecture)
        ├── service.go # Main Service struct and constructor
        ├── get_health.go # Health check handler method
        ├── get_server_info.go # MinIO server info handler method
        ├── get_root.go # Frontend/static asset handler method
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
- **UI Framework**: TailwindCSS 4.1.11 for utility-first styling
- **Design System**: Modern/minimal dashboard theme with gray color palette
- **Theming**: Full dark mode support using `dark:` utility classes

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
# Development mode (default, no embedded assets required)
go run ./main.go -dev

# Production mode with embedded assets
go run -tags dist ./main.go

# Custom port (development mode)
go run ./main.go -dev -addr :3000

# Build binary for verification (development mode, outputs to bin/)
go build -o bin/minio-lite-admin .
./bin/minio-lite-admin -dev

# Build binary with embedded assets (production mode)
go build -tags dist -o bin/minio-lite-admin .
./bin/minio-lite-admin
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

# Pull latest image from GitHub Container Registry
docker pull ghcr.io/[username]/minio-lite-admin:latest
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

### Build Tags System

The project uses Go build tags to conditionally compile embedded assets:

**Development Mode (Default)**:
- File: `embed_dev.go` with `//go:build !dist` tag
- Contains empty `embed.FS` variable
- No `dist/` directory required for compilation
- Command: `go run ./main.go -dev`

**Production Mode**:
- File: `embed_dist.go` with `//go:build dist` tag  
- Contains `//go:embed all:dist` directive
- Requires `dist/` directory with built frontend assets
- Command: `go run -tags dist ./main.go`

This approach prevents build failures when the `dist/` directory doesn't exist during development.

### Build Configuration

**Vite Configuration** (`vite.config.ts`):
```typescript
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  build: {
    manifest: true,        // Required for Go integration
    rollupOptions: {
      input: '/src/main.ts' // Entry point for build
    }
  }
})
```

### TailwindCSS Configuration

**CSS Setup** (`src/style.css`):
```css
@import "tailwindcss";
```

**Key Features**:
- **Version**: TailwindCSS 4.1.11 with official Vite plugin
- **Integration**: `@tailwindcss/vite` plugin for seamless Vite integration
- **Theme**: Modern/minimal design system with professional styling
- **Dark Mode**: Full dark mode support using `dark:` utility classes
- **Responsive**: Mobile-first responsive design with standard breakpoints

**Design Tokens**: See `DESIGN_TOKENS.md` for comprehensive styling guidelines and design system documentation.

### Service Architecture Pattern

The HTTP layer uses a Service-based architecture for better dependency management and testability:

**Service Structure**:
```go
type Service struct {
    config               *config.Config
    logger               zerolog.Logger
    getServerInfoService *service.GetServerInfoService
    distFS               embed.FS
}

func NewService(...deps) (http.Handler, error) {
    // Constructor with dependency injection
    return router, nil
}
```

**Handler Organization**:
- `service.go` - Main Service struct and router configuration
- `get_health.go` - `GetHealthHandler()` method for health checks
- `get_server_info.go` - `GetServerInfoHandler()` method for MinIO info
- `get_root.go` - `GetRootHandler()` method for frontend assets
- `middleware.go` - Shared middleware functions

**Benefits**:
- **Dependency Injection**: All dependencies managed through Service struct
- **Testability**: Easy to mock Service methods for unit testing
- **Organization**: Clear file-to-method mapping for maintainability
- **Consistency**: Unified approach to handler management

### GitHub Actions CI/CD

**Workflow File**: `.github/workflows/docker-publish.yml`

**Features**:
- **Triggers**: Push to `main`/`release` branches, tags (`v*`), and pull requests
- **Multi-platform**: Builds for `linux/amd64` and `linux/arm64` architectures
- **Registry**: Publishes to GitHub Container Registry (GHCR) at `ghcr.io/[owner]/[repo]`
- **Caching**: Uses GitHub Actions cache for faster builds
- **Security**: Includes build attestation for supply chain security
- **Tagging Strategy**:
  - Branch names for branch pushes
  - Tag names for tag pushes  
  - `latest` for main branch
  - Git SHA with branch prefix for main branch
  - PR numbers for pull requests

**Required Permissions**:
- `contents: read` - Access repository code
- `packages: write` - Push to GitHub Container Registry
- `attestations: write` - Generate build attestations
- `id-token: write` - OIDC token for attestations

**Usage**:
```bash
# Pull and run the latest image
docker pull ghcr.io/[owner]/minio-lite-admin:latest
docker run -p 8080:8080 ghcr.io/[owner]/minio-lite-admin:latest
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
- ✅ `/api/server-info` endpoint returns MinIO server information
- ✅ Service-based HTTP handler architecture with dependency injection
- ✅ Modular UI components following clean architecture principles
- ✅ Conditional asset embedding with Go build tags
- ✅ Vue.js 3 + TypeScript frontend scaffold
- ✅ TailwindCSS 4.1.11 integration with modern/minimal dashboard UI
- ✅ Dashboard with server status card, loading states, and responsive design
- ✅ Dark mode support and consistent design tokens
- ✅ Docker development environment with watch mode and MinIO service
- ✅ Production Docker build with multi-stage process
- ✅ Development/production mode switching
- ✅ GitHub Actions workflow for automated Docker builds and GHCR publishing
- ✅ Graceful HTTP server shutdown with signal handling and timeout
- ✅ Comprehensive API testing with table-driven tests and httptest
- ✅ Mock MinIO server infrastructure with chi router for external dependency testing
- ✅ Server-info endpoint testing with success/error scenarios and timeout resolution
- ✅ Code quality tools integration (go test, golangci-lint, gofmt)
- ✅ Disk usage monitoring API with optimized single MinIO API call
- ✅ `/api/data-usage` endpoint returns comprehensive disk metrics
- ✅ Frontend disk usage dashboard with real-time data integration
- ✅ Vue.js composables for API integration with loading and error states
- ✅ TailwindCSS disk usage components with progress bars and status indicators
- ✅ Responsive design with dark mode support for disk usage display
- ✅ Data formatting utilities for human-readable storage values

**Next Steps (TODO)**:
- Additional MinIO administrative features (disk status, access keys, replication)
- Frontend components for MinIO management UI
- MinIO connection health check and validation
- Authentication and authorization system
- More API endpoints for MinIO operations (buckets, users, policies)

## License Constraints

This project uses AGPLv3 license due to dependency on `github.com/minio/madmin-go` which is AGPLv3 licensed. Any derivative work must maintain AGPLv3 compatibility.

## Working Style & Architecture Decisions

This section documents key insights learned from collaborating with Aotokitsuruya and architectural decisions made during development.

### Code Organization Philosophy

**Clean Architecture Principles**:
- Clear separation between domain logic (`internal/service/`) and HTTP concerns (`internal/handler/http/`)
- Dependency injection pattern using struct-based services
- Single responsibility principle with dedicated files per handler
- Interface-based abstractions for testability

**File Naming Conventions**:
- Handler files follow `get_*.go` pattern for HTTP GET endpoints
- Handler methods align with file names: `get_health.go` → `GetHealthHandler()`
- Consistent naming improves code discoverability and maintenance

### Service Architecture Pattern

**Evolution from Functional to Struct-Based Handlers**:
```go
// Old approach (functional)
func HealthHandler(w http.ResponseWriter, r *http.Request) { ... }

// New approach (struct methods)
type Service struct { dependencies... }
func (s *Service) GetHealthHandler(w http.ResponseWriter, r *http.Request) { ... }
```

**Benefits Achieved**:
- Better dependency management and injection
- Improved testability with mockable Service struct
- Cleaner separation of concerns
- More maintainable object-oriented design

### Build System Design

**Conditional Asset Embedding**:
- Development mode as default (no build tags required)
- Production builds use `-tags dist` for embedded assets
- Prevents build failures when `dist/` directory doesn't exist
- Clean separation between development and production concerns

**Key Files**:
- `embed_dev.go` - Default development mode with empty embed.FS
- `embed_dist.go` - Production mode with `//go:embed all:dist`
- Build tags prevent conflicts and compilation errors

### Frontend Integration Strategy

**Hybrid Approach**:
- Development: Vite dev server proxy for hot reloading
- Production: Embedded static assets in Go binary
- Single binary deployment with all assets included

**Critical Design Decision**: 
The `ViteURL` in configuration must be browser-accessible, not just server-accessible. This ensures proper asset loading in containerized environments.

### Testing Strategy

The project follows Go testing best practices with comprehensive API testing:

**Testing Architecture**:
- **Pure Go Testing**: Uses standard `testing.T` without third-party frameworks
- **Table-Driven Tests**: Structured test cases for comprehensive coverage
- **HTTP Testing**: `httptest.ResponseRecorder` for API endpoint testing
- **Dependency Injection**: `testService()` helper for clean test setup

**Test Structure**:
```go
// Test helper in service_test.go
func testService() *Service {
    // Setup minimal dependencies for testing
}

// Table-driven tests in *_test.go files
tests := []struct {
    name           string
    method         string
    expectedStatus int
    expectedBody   ResponseType
}{...}
```

**Testing Patterns**:
- `get_health_test.go` - Health endpoint testing with multiple HTTP methods
- JSON response validation with structure checking
- HTTP header validation (Content-Type, status codes)
- Response format verification (field types, required fields)
- Coverage reporting with `go test -cover`

**Test File Organization**:
- `service_test.go` - Shared test utilities and Service setup
- `get_*_test.go` - Individual handler tests matching file naming convention
- Test files colocated with implementation for easy discovery

**Mock Infrastructure for External Dependencies**:
- **Problem**: Testing endpoints that depend on external services (MinIO Admin API)
- **Solution**: Create mock servers using `httptest.Server` with chi router
- **Location**: `internal/testability/minio/` - Dedicated package for test infrastructure
- **Pattern**: Mock servers simulate external API behavior without network dependencies

**Mock MinIO Server Architecture**:
```go
// Chi router for clean API endpoint management
r := chi.NewRouter()
r.Route("/minio/admin", func(r chi.Router) {
    r.Get("/v4/info", mock.handleServerInfo)
    r.Get("/v3/info", mock.handleServerInfo) 
    r.Get("/info", mock.handleServerInfo)
})

// Configurable responses for different test scenarios
mock.SetServerInfoResponse(scenarios.SuccessfulServerInfo())
mock.SetServerInfoError(400, "Bad Request") // Non-retryable for fast tests
```

### Verification and Quality Practices

**Code Quality Workflow**:
1. Use `gofmt -w .` for formatting verification before builds
2. Run `go test ./...` to execute all tests with coverage reporting
3. Use `golangci-lint` for comprehensive static code analysis
4. Test both development and production build modes
5. Verify functionality with actual HTTP requests
6. Commit only after successful verification and linting

**Troubleshooting Test Timeouts**:
When encountering test timeouts, follow this systematic debugging approach:

1. **Create Simple Test Cases**: Build minimal HTTP servers to isolate the issue
2. **Test Direct HTTP**: Verify mock server responds quickly to direct HTTP requests
3. **Check Concurrent Access**: Test multiple simultaneous requests for lock/blocking issues
4. **Examine Client Behavior**: Investigate if external clients have retry mechanisms
5. **Status Code Selection**: Choose appropriate HTTP status codes (non-retryable for fast tests)

**Common Timeout Causes**:
- **Client Retries**: External libraries may retry on certain HTTP status codes (503, 502, 408, 429)
- **Lock Contention**: Shared resources causing blocking (rarely the actual cause)
- **Network Issues**: DNS resolution or connection timeouts (use httptest to avoid)
- **Infinite Loops**: Logic errors in mock server handlers

**Quality Tools**:
- **go test**: Standard Go testing with coverage analysis
- **golangci-lint**: Comprehensive linting with multiple analyzers
- **gofmt**: Code formatting and syntax verification
- **httptest**: HTTP handler testing without external dependencies

**Development Workflow**:
1. Terminal 1: `pnpm dev` (Vite dev server)
2. Terminal 2: `go run ./main.go -dev` (Go server in dev mode)
3. Browser: http://localhost:8080

### Lessons Learned

**Effective Collaboration Patterns**:
- Start with clear requirements and desired end state
- Break complex refactoring into manageable steps with clear todos
- Verify each step with actual builds and tests
- Use consistent naming conventions for better code organization
- Prefer explicit over implicit (build tags, naming patterns)

**Technical Insights**:
- Go build tags are powerful for conditional compilation
- Struct-based HTTP handlers improve dependency management
- Consistent file/method naming reduces cognitive load
- Service pattern enables better testing and mocking
- Single binary deployment simplifies operations
- Table-driven tests provide comprehensive coverage with minimal code
- httptest.ResponseRecorder enables testing without external dependencies
- Pure Go testing is sufficient for most API testing needs
- Test helpers (testService) reduce duplication and improve maintainability
- golangci-lint catches issues that go compiler misses

**Mock Testing and Timeout Debugging**:
- **External client retry behavior** can cause test timeouts (madmin retries HTTP 503 errors)
- **Systematic debugging approach**: Create simple HTTP servers to isolate issues
- **Root cause identification**: Test timeouts usually stem from retries, not locks/blocking
- **Status code selection matters**: Use non-retryable HTTP codes (400, 401, 404) for fast error tests
- **Chi router in mocks**: Simplifies endpoint management vs hardcoded string matching
- **httptest.Server performance**: Mock servers respond in microseconds, proving no infrastructure issues

**Architecture Evolution**:
The project evolved from simple functional handlers to a well-structured service-based architecture while maintaining backward compatibility and improving maintainability.

### Frontend Development Workflow

**Progressive Enhancement Pattern**:
- Start with placeholder UI components showing static content
- Implement composables for API integration with proper error handling
- Replace static content with real data from backend APIs
- Add loading states, error messages, and user feedback
- Enhance with responsive design and dark mode support

**Vue.js Composition API Patterns**:
```typescript
// Composable pattern for API integration
export function useDataUsage() {
  const dataUsage = ref<DataUsage | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchDataUsage = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('/api/data-usage')
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      dataUsage.value = await response.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  return { dataUsage, loading, error, fetchDataUsage }
}
```

**Component Architecture Principles**:
- **Single Responsibility**: Each component has one clear purpose
- **Composability**: Reusable components that can be combined
- **Props Interface**: Clear TypeScript interfaces for component props
- **Reactive State**: Vue 3 reactivity for dynamic data updates
- **Error Boundaries**: Graceful error handling with fallback UI

**TailwindCSS Integration Patterns**:
```vue
<!-- Responsive design with dark mode support -->
<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
  <div class="bg-white dark:bg-gray-800 rounded-lg p-6">
    <div class="text-gray-900 dark:text-white">
      <!-- Content with proper contrast -->
    </div>
  </div>
</div>
```

**Data Formatting and Utilities**:
- **Utility Functions**: Centralized formatting for bytes, percentages, numbers
- **Type Safety**: TypeScript interfaces for all API responses
- **Computed Properties**: Reactive data transformations in Vue components
- **Internationalization Ready**: Number formatting with locale support

**API Integration Strategy**:
- **Composable Pattern**: Reusable API logic in composables
- **Error Handling**: Comprehensive error states and user feedback
- **Loading States**: Spinner components and skeleton loading
- **Data Transformation**: Backend data adapted to frontend needs
- **Caching Strategy**: Reactive data stores for efficient updates

**Development Workflow Insights**:
1. **Frontend-Backend Coordination**: Ensure API contracts match frontend expectations
2. **Incremental Development**: Build features step by step with frequent commits
3. **Real-Time Integration**: Connect frontend to live backend APIs early
4. **Component Testing**: Verify UI components with actual data flows
5. **Responsive Design**: Test across different screen sizes and themes

**Vue.js Best Practices Learned**:
- **Composition API**: Superior to Options API for complex logic
- **`<script setup>`**: Cleaner syntax for component setup
- **Computed Properties**: For reactive data transformations
- **Watchers**: For side effects and API calls on data changes
- **Lifecycle Hooks**: Proper component initialization with `onMounted`

**Frontend Architecture Benefits**:
- **Modularity**: Composables enable code reuse across components
- **Maintainability**: Clear separation between UI and business logic
- **Testability**: Pure functions in utilities and composables
- **Performance**: Vue 3 reactivity system optimizes re-renders
- **Developer Experience**: TypeScript and Vite provide excellent DX

**Integration with Backend Services**:
- **API Endpoints**: RESTful endpoints with consistent JSON responses
- **Error Standardization**: Backend errors properly handled in frontend
- **Data Contracts**: TypeScript interfaces ensure type safety
- **Development Mode**: Hot reloading with backend changes
- **Production Builds**: Optimized bundles embedded in Go binary

## Development Notes

- The project requires Go 1.24+ and Node.js for development
- Frontend assets must be built (`pnpm build`) before production deployment
- The `dist/.keep` file ensures the dist directory exists for Go embed
- API endpoints are structured under `/api` prefix for clear separation
- Static assets are served by Vite dev server in development, embedded in Go binary in production
- Build output goes to `bin/` directory (excluded from git) for verification without cleanup
- Use `go run ./main.go -dev` for development, `go build . -o bin/minio-lite-admin` for verification
- Use `gofmt -w .` to format all Go files before running verification tests
- Run `go test ./...` to execute all tests
- Run `go test ./... -cover` to execute all tests with coverage analysis
- Use `golangci-lint run` for comprehensive static code analysis and linting
- Test specific endpoints: `go test -v ./internal/handler/http -run TestService_GetServerInfoHandler`
- Run tests with timeout limits: `go test -v ./internal/handler/http -timeout 30s`
- Generate detailed coverage report: `go test ./internal/handler/http -coverprofile=coverage.out && go tool cover -func=coverage.out`