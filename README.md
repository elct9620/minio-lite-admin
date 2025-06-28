# MinIO Lite Admin

A lightweight web-based administration tool for MinIO instances, built with Go and Vue.js.

## Overview

MinIO Lite Admin provides essential management operations for MinIO instances through a user-friendly web interface. Following MinIO's decision to remove most management features from the community edition, this tool fills the gap by providing access to commonly used administrative functions.

## Features

- **🖥️ Web UI Dashboard** - Modern Vue.js interface with dark mode support
- **💾 Disk Usage Monitoring** - Real-time disk status and usage statistics  
- **🔑 Access Key Management** - Create, list, update, and delete access keys with secure generation
- **📊 Server Information** - View MinIO server status and configuration
- **🔒 Secure Key Generation** - Cryptographically secure access key generation with special characters

## Quick Start

### Docker (Production)

Run MinIO Lite Admin as a standalone container:

```bash
docker run -p 8080:8080 \
  -e MINIO_URL=http://your-minio-server:9000 \
  -e MINIO_ROOT_USER=your-admin-user \
  -e MINIO_ROOT_PASSWORD=your-admin-password \
  ghcr.io/[owner]/minio-lite-admin:latest
```

### Docker Compose (Side-car with MinIO)

Create a `docker-compose.yml` file to run alongside your MinIO instance:

```yaml
services:
  # Your existing MinIO service
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  # MinIO Lite Admin (side-car)
  minio-admin:
    image: ghcr.io/[owner]/minio-lite-admin:latest
    ports:
      - "8080:8080"
    environment:
      - MINIO_URL=http://minio:9000
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    depends_on:
      - minio

volumes:
  minio_data:
```

Then run:

```bash
docker compose up -d
```

Access the admin interface at: http://localhost:8080

## Configuration

Configure MinIO Lite Admin using environment variables:

### Required Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `MINIO_URL` | MinIO server URL | `http://minio:9000` |
| `MINIO_ROOT_USER` | MinIO admin username | `minioadmin` |
| `MINIO_ROOT_PASSWORD` | MinIO admin password | `minioadmin` |

### Optional Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MINIO_ADMIN_ADDR` | `:8080` | Server bind address |
| `MINIO_ADMIN_LOG_LEVEL` | `info` | Log level (trace, debug, info, warn, error) |
| `MINIO_ADMIN_LOG_PRETTY` | `true` | Pretty print logs |

### Development Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MINIO_ADMIN_DEV` | `false` | Enable development mode |
| `MINIO_ADMIN_VITE_URL` | `http://localhost:5173` | Vite dev server URL |

## Security

### Access Key Generation

- **Cryptographically Secure**: Uses `crypto.getRandomValues()` for secure random generation
- **Enhanced Character Set**: Supports special characters for higher entropy
- **AWS Compatible**: Follows AWS IAM access key format with `AKIA` prefix
- **MinIO Compatible**: All generated keys work seamlessly with MinIO backend

### Best Practices

- Always use strong MinIO admin credentials
- Run behind HTTPS in production
- Restrict network access to admin interface
- Regularly rotate access keys

## Development

For development setup and contribution guidelines, see [CLAUDE.md](CLAUDE.md).

### Quick Development Setup

```bash
# Clone repository
git clone <repository-url>
cd minio-lite-admin

# Start development environment
docker compose up --watch
```

This starts:
- MinIO server at http://localhost:9000
- Vite dev server at http://localhost:5173  
- Go backend at http://localhost:8080

## Tech Stack

- **Backend**: Go 1.24+
  - MinIO Admin SDK (`github.com/minio/madmin-go`)
  - Chi HTTP router with structured logging
  - Vite integration for seamless asset serving
- **Frontend**: Vue.js 3 + TypeScript
  - Composition API with `<script setup>`
  - TailwindCSS for styling with dark mode
  - Vite for fast development and optimized builds

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3) - see the [LICENSE](LICENSE) file for details.

### Why AGPLv3?

This project uses `github.com/minio/madmin-go` which is licensed under AGPLv3. As a derivative work, this project must also be licensed under AGPLv3 to comply with the license terms.

## Support

- 📖 Documentation: [CLAUDE.md](CLAUDE.md)
- 🐛 Issues: [GitHub Issues](https://github.com/[owner]/minio-lite-admin/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/[owner]/minio-lite-admin/discussions)

## Acknowledgments

- MinIO team for the excellent `madmin-go` SDK
- Vue.js and TailwindCSS communities for amazing tools
- All contributors who help improve this project