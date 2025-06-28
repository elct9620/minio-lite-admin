# Multi-stage production Dockerfile for MinIO Lite Admin

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package files and install dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Copy frontend source and configuration
COPY vite.config.ts tsconfig*.json ./
COPY src/ ./src/
COPY public/ ./public/

# Build frontend for production
RUN pnpm build

# Stage 2: Build backend with embedded frontend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY . .

# Copy built frontend assets from previous stage
COPY --from=frontend-builder /app/dist ./dist/

# Build Go binary with embedded assets
RUN CGO_ENABLED=0 GOOS=linux go build -tags dist -o main .

# Stage 3: Final production image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=backend-builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main", "-addr", ":8080"]