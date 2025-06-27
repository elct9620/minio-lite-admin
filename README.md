# MinIO Lite Admin

A lightweight web-based administration tool for MinIO instances, built with Go and Vue.js.

## Overview

MinIO Lite Admin is designed to provide essential management operations for MinIO instances through a user-friendly web interface. This project utilizes the official `github.com/minio/madmin-go` SDK to interact with MinIO servers.

### Background

Following MinIO's decision in May 2025 to remove most management features from the community edition, this tool aims to fill the gap by providing access to commonly used administrative functions. While not intended to be a complete replacement for all MinIO management features, it focuses on the most frequently needed operations.

## Features

The project focuses on essential MinIO management tasks:

- **Disk Status Monitoring** - Check and monitor disk health and usage
- **Access Key Management** - Create, list, and manage access keys
- **Site Replication Configuration** - Configure and manage site replication settings

## Tech Stack

- **Backend**: Go 1.24 (latest)
  - MinIO Admin SDK (`github.com/minio/madmin-go`)
  - Vite integration via `github.com/olivere/vite`
- **Frontend**: Vue.js
  - Vite as the build tool
  - Modern reactive UI for better user experience

## Prerequisites

- Go 1.24 or higher
- Node.js and npm (for frontend development)
- Access to a MinIO instance with admin credentials

## Installation

TODO

## Configuration

TODO

## Usage

TODO

## Development

### Backend Development

TODO

### Frontend Development

TODO

### Building from Source

TODO

## API Documentation

TODO

## Contributing

TODO

## Security

TODO

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3) - see the [LICENSE](LICENSE) file for details.

### Why AGPLv3?

This project uses `github.com/minio/madmin-go` which is licensed under AGPLv3. As a derivative work, this project must also be licensed under AGPLv3 to comply with the license terms.

## Support

TODO

## Roadmap

TODO

## Acknowledgments

- MinIO team for the excellent `madmin-go` SDK
- All contributors who help improve this project