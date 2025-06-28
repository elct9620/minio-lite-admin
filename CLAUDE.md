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
go run . -dev

# Production mode with embedded assets
go run -tags dist .

# Custom port (development mode)
go run . -dev -addr :3000

# Build binary for verification (development mode, outputs to bin/)
go build -o bin/minio-lite-admin .
./bin/minio-lite-admin -dev

# Build binary with embedded assets (production mode)
go build -tags dist -o bin/minio-lite-admin .
./bin/minio-lite-admin
```

### Full Development Workflow
1. Terminal 1: `pnpm dev` (starts Vite dev server)
2. Terminal 2: `go run . -dev` (starts Go server in dev mode)
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

### Quality Assurance
```bash
# Format code (required before commits)
gofmt -w .

# Run all tests with race detection and coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate and view coverage report
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run comprehensive linting (matches CI pipeline)
golangci-lint run --timeout=5m

# Test specific package
go test -v ./internal/handler/http

# Test specific function
go test -v ./internal/handler/http -run TestService_GetServerInfoHandler

# Run tests with timeout
go test -v ./internal/handler/http -timeout 30s

# Clean up test artifacts
rm coverage.out coverage.html
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
- Command: `go run . -dev`

**Production Mode**:
- File: `embed_dist.go` with `//go:build dist` tag  
- Contains `//go:embed all:dist` directive
- Requires `dist/` directory with built frontend assets
- Command: `go run -tags dist .`

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

**Workflow File**: `.github/workflows/ci.yml`

**Pipeline Architecture**:
- **Test Job**: Go 1.24 with race detection, coverage reporting, and module caching
- **Lint Job**: Code formatting validation (gofmt) and comprehensive linting (golangci-lint)
- **Build-and-Push Job**: Docker image building and publishing (depends on test and lint passing)

**Features**:
- **Triggers**: Push to `main`/`release` branches, tags (`v*`), and pull requests
- **Quality Gates**: All tests and linting must pass before Docker builds
- **Multi-platform**: Builds for `linux/amd64` and `linux/arm64` architectures
- **Registry**: Publishes to GitHub Container Registry (GHCR) at `ghcr.io/[owner]/[repo]`
- **Caching**: Uses GitHub Actions cache for Go modules and Docker layers
- **Security**: Includes build attestation for supply chain security
- **Code Quality**: Automated gofmt checking and golangci-lint validation
- **Test Coverage**: Comprehensive test execution with race detection
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

### Current Implementation Status

**Core Features Implemented**:
- ✅ `/api/server-info` - MinIO server information
- ✅ `/api/data-usage` - Disk usage monitoring  
- ✅ `/api/access-keys` - Access key management (GET, POST, PUT, DELETE)
- ✅ Vue.js frontend with server dashboard, disk usage, and access key management
- ✅ Docker development environment with MinIO service

**Architecture Components**:
- ✅ Go backend with Chi router and zerolog middleware
- ✅ Service layer architecture (`internal/service/`) for business logic
- ✅ Mock MinIO server for testing (`internal/testability/minio/`)
- ✅ Build tags system for dev/production asset embedding
- ✅ Vue.js 3 + TypeScript frontend with TailwindCSS

**Development Infrastructure**:
- ✅ CI/CD pipeline with Go Test, Go Lint, and Docker build
- ✅ Release-please for automated version management
- ✅ Table-driven testing patterns with comprehensive coverage
- ✅ Docker development with hot reload support

**Next Major Features**:
- Site replication configuration (placeholder view exists)
- Authentication/authorization system
- Additional MinIO admin operations

## License Constraints

This project uses AGPLv3 license due to dependency on `github.com/minio/madmin-go` which is AGPLv3 licensed. Any derivative work must maintain AGPLv3 compatibility.

## Architecture Decisions

### Service Architecture Pattern
- **HTTP Handlers**: Struct-based handlers with dependency injection via `Service` struct
- **File Naming**: Handler files follow `get_*.go` pattern with matching method names
- **Testing**: Each handler has corresponding `*_test.go` with table-driven tests

### Build System Design
- **Development Mode**: Default mode, no build tags required, uses Vite dev server
- **Production Mode**: `-tags dist` embeds frontend assets in Go binary
- **Critical**: `ViteURL` must be browser-accessible for containerized environments

### Frontend Integration
- **Development**: Go server proxies to Vite dev server (localhost:5173)
- **Production**: Go serves embedded assets from `dist/` directory
- **Template System**: HTML templates with Vite fragment injection (`{{ .Vite.Tags }}`)

### Testing Strategy
- **Table-Driven Tests**: Standard pattern with `httptest.ResponseRecorder`
- **Mock Infrastructure**: `internal/testability/minio/` provides mock MinIO server using chi router
- **Test Naming**: `get_*_test.go` files match handler files with `testService()` helper
- **External Dependencies**: Mock servers simulate MinIO Admin API with proper v4 encryption

### MinIO Admin API Integration
- **Access Keys API**: Unified interface for users, service accounts, and STS keys via `ListAccessKeysBulk()`
- **Service Account Creation**: `AddServiceAccount()` auto-generates keys when not provided, accepts custom keys
- **madmin-go v4**: Requires encrypted responses in mock servers using `madmin.EncryptData()`
- **API Options**: Mutually exclusive options require conditional handling (`All: true` vs specific users)
- **Constants**: Always use madmin package constants (e.g., `madmin.AccessKeyListUsersOnly`) for type safety
- **Policy Configuration**: JSON policies passed as `[]byte` in `AddServiceAccountReq.Policy`

### API Response Format and Frontend Integration

**Critical Issue: JSON Array Handling**
- **Problem**: Go `nil` slices serialize to `null`, causing frontend iteration errors
- **Solution**: Always initialize slices to `[]` before JSON marshaling
- **Detection**: Test API endpoints with `curl` before debugging frontend

**Vue 3 Integration Patterns**
- **Composables**: Reusable API logic with reactive refs for loading, error, and data states
- **Component Architecture**: `<script setup>` with TypeScript, proper loading/error/empty states
- **SPA Routing**: Backend serves `index.html` for non-API routes using `isStaticAsset()` detection
- **Icon System**: Use Heroicons (`@heroicons/vue`) with dynamic component rendering for consistency
- **Modal Components**: Use overlay patterns with backdrop click handling and form state management

### MinIO Admin API Method Research

**Essential Resource**: Use `go doc github.com/minio/madmin-go/v4.AdminClient.MethodName` to understand:
- **Endpoint Paths**: All admin methods use `/minio/admin/v4/endpoint-name` pattern
- **HTTP Methods**: GET, POST, PUT, DELETE based on operation type
- **Parameters**: Query parameters vs request body structure
- **Response Format**: Expected status codes and data structures

**Key Examples**:
- `DeleteServiceAccount`: DELETE `/v4/delete-service-account?accessKey=KEY` → HTTP 204
- `AddServiceAccount`: PUT `/v4/add-service-account` + encrypted body → HTTP 200 + encrypted response
- `ListAccessKeysBulk`: GET `/v4/list-access-keys-bulk` + query params → HTTP 200 + encrypted response

**Debug Process**:
1. Use `go doc` to check madmin method signature and behavior
2. Check mock server endpoint registration matches actual MinIO paths
3. Verify query parameter names match exactly (e.g., `accessKey` not `service-account`)
4. Ensure encrypted request/response handling for methods that require it

### Frontend Modal Pattern

**Component Structure** for delete/confirmation workflows:
- **Props**: `open: boolean`, specific data fields (accessKey, name, etc.)
- **Emits**: `'update:open'` for v-model, `'confirmed'` for action completion
- **State Management**: `loading`, `error`, API call handling within modal
- **UI Elements**: Warning icon, danger styling, clear action buttons
- **Error Handling**: Display API errors inline within modal context

**Integration Pattern**:
- Parent manages modal state (`showModal`, `itemToDelete`)
- Event handlers: open modal → API call → close + refresh
- Conditional rendering: `v-if="itemToDelete"` to ensure props are available

## Security Implementation

### Access Key Generation (Frontend)
- **Location**: `CreateAccessKeyModal.vue:159-179`, `EditAccessKeyModal.vue:226-237`
- **Method**: `crypto.getRandomValues()` for cryptographically secure generation
- **Access Key Format**: `AKIA` prefix + 16 secure chars (AWS-compatible)
- **Secret Key Charset**: Alphanumeric + special chars `+/=!@#$%^&*()_-[]{}|;:,.<>?` (79 chars total)
- **Length**: Access keys = 20 chars, Secret keys = 40 chars
- **Security**: High entropy, cryptographically secure, MinIO compatible

### Key Security Patterns
- Always use `crypto.getRandomValues()` over `Math.random()` for credentials
- Backend accepts any custom keys without validation (rely on MinIO validation)
- MinIO auto-generates secure keys when `AccessKey`/`SecretKey` are empty
- Special characters in secret keys are fully supported by MinIO backend

## Coding Style

### Go Conventions
- **Type Alias**: Use `any` instead of `interface{}` for better readability (Go 1.18+ standard)
- **Error Handling**: Always handle errors explicitly, avoid ignoring with `_`
- **Naming**: Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- **Package Names**: Short, lowercase, single words when possible

## Development Notes

**Requirements**: Go 1.24+ and Node.js

**Key Patterns**:
- API endpoints under `/api` prefix, frontend routes served via SPA fallback
- Development: Vite dev server, Production: embedded assets with `-tags dist`
- **Build Commands**: Use `go run .` (not `go run ./main.go`) to include embed files correctly
- All slice responses must be initialized to `[]` (not `nil`) for JSON compatibility
- Mock servers require madmin v4 encryption for MinIO Admin API testing
- Service creation: MinIO auto-generates secure keys when `AccessKey`/`SecretKey` are empty
- Form validation: Always validate required fields before API calls, show clear error messages

## Release Management

### Release-Please Integration
- **Automated Releases**: Uses `googleapis/release-please-action@v4` for semantic versioning
- **Configuration**: Simple workflow in `.github/workflows/release-please.yml`
- **Permissions**: Requires `contents: write`, `issues: write`, `pull-requests: write`
- **Initial Version**: Set using `Release-As: 0.1.0` footer in commit message
- **Version Strategy**: Conventional commits drive version bumps (fix → patch, feat → minor, feat! → major)

### Conventional Commits
- **Required**: All commits must follow conventional commit format for release automation
- **Examples**: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`
- **Breaking Changes**: Use `!` suffix (e.g., `feat!:`) for major version bumps
- **AI Attribution**: Include `Co-Authored-By: Claude <noreply@anthropic.com>` for AI-assisted development

## Project Maintenance

### AI-Assisted Development
- **Policy**: AI tools are welcome, as stated in README.md Contributing section
- **Attribution**: Include `Co-Authored-By: Claude <noreply@anthropic.com>` for AI-assisted commits
- **PR Template**: Simple template in `.github/pull_request_template.md` requires AI usage declaration

### GitHub Workflows
- **CI/CD Pipeline**: `.github/workflows/ci.yml` with Go Test, Go Lint, and Docker build jobs
- **Release Automation**: `.github/workflows/release-please.yml` for automated version management
- **PR Template**: `.github/pull_request_template.md` with AI usage declaration and basic requirements

### Recent Development Insights
- **Dependabot Management**: Use `@dependabot rebase` to fix lint issues after main branch updates
- **Lint Resolution**: Modern CI improvements resolved previous golangci-lint failures
- **Docker Improvements**: Resolved build issues with proper source directory copying
- **Workflow Optimization**: Better job names and streamlined CI process

### Release-Please Usage

**Setup**: Release-please is configured in `.github/workflows/release-please.yml` for automated version management.

**How to Use Release-Please**:
1. **Commit with conventional format**: `feat:`, `fix:`, `docs:`, `chore:`
2. **Push to main**: Release-please analyzes commits and creates release PR when needed
3. **Review release PR**: Check generated CHANGELOG.md and version bump
4. **Merge release PR**: Creates GitHub release and git tag automatically

**Initial Version Setup**:
- Use `Release-As: X.Y.Z` footer in commit message to set specific version
- Example: `git commit --allow-empty -m "chore: initialize release versioning" -m "Release-As: 0.1.0"`

**Version Bumping Rules**:
- `fix:` → patch version (0.1.0 → 0.1.1)
- `feat:` → minor version (0.1.0 → 0.2.0)  
- `feat!:` or `fix!:` → major version (0.1.0 → 1.0.0)

**Current Codebase Status**:
- Server info and data usage APIs implemented
- Access key management (create, list, update, delete) implemented
- Vue.js frontend with TailwindCSS and dark mode
- Docker development environment with MinIO service
- CI/CD pipeline with test, lint, and Docker build
- Release-please configured for version 0.1.0