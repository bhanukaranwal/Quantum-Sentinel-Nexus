# services/alert-router-svc/Dockerfile

# --- Stage 1: Builder ---
# This stage installs all development tools and libraries, then compiles the application.
FROM debian:bookworm as builder

# Set environment variables to avoid interactive prompts during installation.
ENV DEBIAN_FRONTEND=noninteractive

# Install build tools and dependencies for our C++ application.
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    pkg-config \
    librdkafka-dev \
    libzmq3-dev \
    libssl-dev \
    libuv1-dev \
    zlib1g-dev

# Set the working directory.
WORKDIR /app

# --- Dependency: uWebSockets ---
# uWebSockets is a key dependency that we'll build from source.
RUN git clone https://github.com/uNetworking/uWebSockets.git
WORKDIR /app/uWebSockets
RUN git checkout v20.40.0 && make && make install

# --- Dependency: nlohmann/json ---
# This is a header-only library, so we just need to download it.
WORKDIR /app
RUN git clone https://github.com/nlohmann/json.git vendor/json

# --- Build the Application ---
# Copy the application source code into the container.
COPY . .

# Create a build directory and run CMake and make.
RUN cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
RUN cmake --build build

# --- Stage 2: Final Image ---
# This stage creates the final, minimal image for deployment.
FROM debian:bookworm-slim

# Install only the runtime dependencies, not the -dev packages.
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    libzmq5 \
    libssl3 \
    libuv1 \
    zlib1g \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app
USER app

# Set the working directory.
WORKDIR /app

# Copy the compiled binary from the 'builder' stage.
COPY --from=builder /app/build/alert-router-svc .

# Expose the WebSocket port and the ZeroMQ port.
EXPOSE 8081
EXPOSE 5555

# The command to run when the container starts.
CMD ["./alert-router-svc"]
