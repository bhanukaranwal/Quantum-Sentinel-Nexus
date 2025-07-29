# services/edge-node-agent/Dockerfile

# --- Stage 1: Builder ---
# This stage uses an NVIDIA CUDA development image as its base. This provides
# the CUDA toolkit, compiler (nvcc), and necessary drivers for building the application.
FROM nvidia/cuda:12.1.1-devel-ubuntu22.04 as builder

# Set environment variables to avoid interactive prompts during installation.
ENV DEBIAN_FRONTEND=noninteractive

# Install build tools and dependencies for our C++ application.
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    pkg-config \
    libgrpc-dev \
    libgrpc++-dev \
    protobuf-compiler \
    protobuf-compiler-grpc \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory.
WORKDIR /app

# Copy the application source code into the container.
COPY . .

# Create a build directory and run CMake and make.
# This will find the pre-installed CUDA and gRPC libraries and build the agent.
RUN cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
RUN cmake --build build

# --- Stage 2: Final Image ---
# This stage creates the final, minimal image for deployment.
# We use a CUDA 'base' image, which contains the runtime libraries but not the
# development toolkit, making it much smaller.
FROM nvidia/cuda:12.1.1-base-ubuntu22.04

# Install only the runtime dependencies.
RUN apt-get update && apt-get install -y \
    libgrpc++1 \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app
USER app

# Set the working directory.
WORKDIR /app

# Copy the compiled binary from the 'builder' stage.
COPY --from=builder /app/build/edge-node-agent .

# The command to run when the container starts.
# This agent would typically be started with specific configuration pointing
# to its local data and the central coordinator.
CMD ["./edge-node-agent"]
