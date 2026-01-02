# fthenoise

A Go web server for reading and exploring classic texts.

## Features

- Home page with styled HTML
- Health check endpoint (`/health`)
- Hello API endpoint (`/api/hello`)
- Reading interface (`/read`) for exploring classic texts
- Breathe endpoint (`/breathe`)
- JSON responses for API endpoints

## Running with Docker

The recommended way to run the application is using Docker Compose:

```bash
docker-compose up --build
```

This will:

- Build the Docker image using the multi-stage Dockerfile
- Start the container on port 8080
- Automatically restart the container if it stops

The server will be available at `http://localhost:8080`

**Note:** The Docker setup does not support hot reloading. For development with hot reloading, see the [Development](#development) section below.

### Docker Commands

- **Start in detached mode**: `docker-compose up -d`
- **Stop the container**: `docker-compose down`
- **View logs**: `docker-compose logs -f`
- **Rebuild after changes**: `docker-compose up --build`

## Running Locally (Without Docker)

To run the server directly:

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified by the `PORT` environment variable).

## Endpoints

- `GET /` - Home page with server information
- `GET /health` - Health check endpoint (returns JSON)
- `GET /api/hello` - Hello API endpoint (returns JSON)
- `GET /breathe` - Breathe endpoint
- `GET /read` - Reading interface for exploring texts

## Building

To build the server:

```bash
go build -o fthenoise main.go
./fthenoise
```

## Development

### Hot Reloading

This project supports hot reloading when running locally (not in Docker). The project uses [fresh](https://github.com/gravityblast/fresh) for auto-reloading during development. It's tracked as a dev dependency in `tools.go`.

To run with hot reloading:

```bash
fresh
```

This will watch your files and automatically rebuild and restart the server when changes are detected. The server will start on `http://localhost:8080`.

**Note:** Hot reloading only works when running locally. The Docker setup builds a static binary and does not support file watching or hot reloading.

### Installing Dev Dependencies

Dev dependencies are automatically tracked via `tools.go`. To ensure all dependencies are up to date:

```bash
go mod tidy
```

## Project Structure

- `main.go` - Main application code
- `templates/` - HTML templates
- `texts/` - Text files to be loaded and parsed
- `Dockerfile` - Multi-stage Docker build configuration
- `docker-compose.yml` - Docker Compose configuration
