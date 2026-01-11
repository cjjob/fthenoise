# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /build

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -trimpath \
    -o fthenoise main.go

# Runtime stage - use distroless for smallest secure image
FROM gcr.io/distroless/static-debian13:nonroot

# Set working directory
WORKDIR /app

# Copy the binary
COPY --from=builder /build/fthenoise /app/fthenoise

# Copy templates and texts
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/texts /app/texts

# Expose port
EXPOSE 8080

# Use non-root user (distroless provides this)
USER nonroot:nonroot

# Run the binary
ENTRYPOINT ["/app/fthenoise"]

