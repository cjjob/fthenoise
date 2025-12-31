# Go Web Server

A simple, modern Go web server with a beautiful UI.

## Features

- Home page with styled HTML
- Health check endpoint (`/health`)
- Hello API endpoint (`/api/hello`)
- JSON responses for API endpoints

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Endpoints

- `GET /` - Home page with server information
- `GET /health` - Health check endpoint (returns JSON)
- `GET /api/hello` - Hello API endpoint (returns JSON)

## Building

To build the server:

```bash
go build -o server main.go
./server
```

## Development

### Development Dependencies

This project uses [fresh](https://github.com/gravityblast/fresh) for auto-reloading during development. It's tracked as a dev dependency in `tools.go`.

To use fresh for auto-reloading:

```bash
fresh
```

This will watch your files and automatically rebuild and restart the server when changes are detected.

### Installing Dev Dependencies

Dev dependencies are automatically tracked via `tools.go`. To ensure all dependencies are up to date:

```bash
go mod tidy
```
