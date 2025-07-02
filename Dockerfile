# =========================================================================
# Procure App: Full-Stack Dockerfile (Production Ready)
# =========================================================================
# This multi-stage Dockerfile builds a self-contained production image
# for a Go backend and a SvelteKit frontend.
#
# It uses build-time arguments to inject configuration into the frontend assets.
# =========================================================================


# =========================================================================
# STAGE 1: Build Frontend (SvelteKit)
# =========================================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Define the build-time arguments. These names MUST EXACTLY MATCH what is
# imported in your SvelteKit code (e.g., from '$env/static/public').
ARG PUBLIC_VITE_AUTH0_DOMAIN
ARG PUBLIC_VITE_AUTH0_CLIENT_ID
ARG PUBLIC_VITE_AUTH0_CALLBACK_URL
ARG PUBLIC_VITE_AUTH0_AUDIENCE
ARG PUBLIC_VITE_API_BASE_URL

# Copy package management files to leverage Docker's layer cache.
COPY frontend/package*.json ./

# Install dependencies using 'npm ci' for fast, reproducible builds.
RUN npm ci

# Copy the rest of the frontend source code.
COPY frontend/ ./

# Build the frontend. The build-time ARGs are passed as environment variables,
# allowing Vite to bake them into the static assets.
RUN PUBLIC_VITE_AUTH0_DOMAIN=$PUBLIC_VITE_AUTH0_DOMAIN \
    PUBLIC_VITE_AUTH0_CLIENT_ID=$PUBLIC_VITE_AUTH0_CLIENT_ID \
    PUBLIC_VITE_AUTH0_CALLBACK_URL=$PUBLIC_VITE_AUTH0_CALLBACK_URL \
    PUBLIC_VITE_AUTH0_AUDIENCE=$PUBLIC_VITE_AUTH0_AUDIENCE \
    PUBLIC_VITE_API_BASE_URL=$PUBLIC_VITE_API_BASE_URL \
    npm run build


# =========================================================================
# STAGE 2: Build Backend (Go)
# =========================================================================
FROM golang:1.24-bullseye AS backend-builder

WORKDIR /app

# Install build-time dependencies for CGO (for go-sqlite3).
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev libsqlite3-dev && rm -rf /var/lib/apt/lists/*

# Copy and cache Go module dependencies.
COPY backend/go.* ./
RUN go mod download

# Copy the backend source code.
COPY backend/ ./

# Compile the Go application, creating a small, optimized binary.
RUN CGO_ENABLED=1 GOOS=linux go build -mod=mod -ldflags="-s -w" -o /main .


# =========================================================================
# STAGE 3: Final Production Image
# =========================================================================
FROM debian:bullseye-slim

WORKDIR /app

# Install only the required RUNTIME libraries.
RUN apt-get update && apt-get install -y --no-install-recommends sqlite3 ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the built frontend assets from the 'frontend-builder' stage.
COPY --from=frontend-builder /app/frontend/build /app/frontend/dist

# Copy the compiled Go binary from the 'backend-builder' stage.
COPY --from=backend-builder /main /app/main

# --- Security: Run as a non-root user ---
RUN useradd --system --create-home --shell /bin/false appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose the application port.
EXPOSE 8080

# Define the container's startup command.
CMD ["./main"]