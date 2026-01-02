.PHONY: build run stop clean docker-build docker-run docker-stop

# Build the Docker image
docker-build:
	docker build -t fthenoise:latest .

# Run with docker-compose
docker-run:
	docker-compose up -d

# Stop docker-compose
docker-stop:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Build and run
docker-up: docker-build docker-run

# Clean up
clean:
	docker-compose down -v
	docker rmi fthenoise:latest || true

# Local development (without Docker)
build:
	go build -o fthenoise main.go

run:
	go run main.go

# Test the Docker image locally
test-docker: docker-build
	docker run --rm -p 8080:8080 fthenoise:latest

