# services/ingestion-svc/Dockerfile

# --- Build Stage ---
# This stage uses the official Go image to build our application binary.
# We name this stage 'builder' so we can refer to it later.
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files and download dependencies.
# This is done in a separate layer to leverage Docker's build cache.
# Dependencies will only be re-downloaded if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the Go application.
# CGO_ENABLED=0 disables Cgo, which is necessary for creating a static binary.
# -o /ingestion-svc specifies the output path for the compiled binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /ingestion-svc .


# --- Final Stage ---
# This stage creates the final, minimal image.
# We use a "distroless" base image from Chainguard, which contains only our application
# and its direct runtime dependencies. It does not contain a shell or other
# unnecessary programs, which significantly reduces the attack surface.
FROM cgr.dev/chainguard/static:latest

# Copy the compiled binary from the 'builder' stage.
COPY --from=builder /ingestion-svc /ingestion-svc

# Copy any necessary non-code assets, like configuration files or certificates.
# For this service, we don't have any, but this is where they would go.
# COPY --from=builder /app/config.yaml /config/config.yaml

# Expose the port the service listens on.
EXPOSE 8080

# Set the user to 'nonroot' for security. Distroless images have this user by default.
USER nonroot

# The command to run when the container starts.
ENTRYPOINT ["/ingestion-svc"]
