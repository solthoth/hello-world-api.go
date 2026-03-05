# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A minimal Go REST API using [Gin](https://github.com/gin-gonic/gin), designed to run in Kubernetes. The binary is compiled to a static binary and shipped in a `scratch` Docker image.

## Routes

| Method | Path       | Description                  |
|--------|------------|------------------------------|
| GET    | `/`        | Returns `{"message": "Hello, World!"}` |
| GET    | `/health`  | Returns `{"status": "ok"}` — used as a Kubernetes liveness probe |
| POST   | `/message` | Publishes `{"message":"hello world","timestamp":"<utc>"}` to Azure Storage Queue; returns `202` with the same payload |

## Environment Variables

| Variable | Description |
|---|---|
| `AZURE_STORAGE_CONNECTION_STRING` | Azure Storage connection string — sourced from secret `hello-world-api-queue-conn` key `primaryQueueConnectionString` (Crossplane-provisioned) |
| `AZURE_STORAGE_QUEUE_NAME` | Azure queue name — `hello-world-api-queue-queue` (XR name + `-queue` suffix from the composition) |

If either variable is unset the app starts normally but `POST /message` returns `503`.

## Common Commands

```bash
# Run locally
go run .

# Build binary
go build -o hello-world-api .

# Run all tests
go test ./...

# Run a single test by name
go test ./... -run TestName

# Add/update dependencies
go get <package>
go mod tidy
```

## Docker

Multi-stage build: compiles with `golang:1.26-alpine`, runs in `scratch` (no base OS, ~10MB final image).

```bash
docker build -t hello-world-api .
docker run -p 8080:8080 hello-world-api
```

The binary is built with `CGO_ENABLED=0` and `-ldflags="-w -s"` for a static, stripped binary.

## GitHub Actions

The workflow at [.github/workflows/docker-publish.yml](.github/workflows/docker-publish.yml):
- Triggers on push to `main` (builds and pushes) and on pull requests (builds only)
- Publishes to `ghcr.io/solthoth/hello-world-api.go`
- Tags: `latest` (on main) and `sha-<commit>` (always)
- Uses `GITHUB_TOKEN` automatically — no secrets to configure

## Kubernetes Health Check

Point liveness/readiness probes at `GET /health` on port `8080`.
