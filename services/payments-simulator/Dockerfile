# services/payments-simulator/Dockerfile

# --- Stage 1: Go Builder ---
# This stage builds the Go application that will act as a simple
# HTTP server to trigger the Locust tests.
FROM golang:1.22-alpine AS go-builder

WORKDIR /app

# Copy Go module files and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source code.
COPY . .

# Build the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build -o /payments-simulator .


# --- Stage 2: Final Image ---
# This stage creates the final image containing Python, Locust,
# the compiled Go app, and the locustfile.
FROM python:3.12-slim

WORKDIR /app

# Install Locust.
RUN pip install locust==2.28.0

# Copy the compiled Go binary from the 'go-builder' stage.
COPY --from=go-builder /payments-simulator /usr/local/bin/

# Copy the locustfile that defines the user behavior.
COPY locustfile.py .

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app
USER app

# Expose the port for the Go HTTP server (which also serves the Locust UI).
EXPOSE 8089

# The command to run when the container starts.
# This starts the Go application, which in turn can trigger Locust.
CMD ["/usr/local/bin/payments-simulator"]
