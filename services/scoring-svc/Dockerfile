# services/scoring-svc/Dockerfile

# --- Stage 1: Planner ---
# This stage's only purpose is to cache our dependencies.
# It copies only the files needed to resolve the dependency tree
# and builds a dummy project. This layer will only be rebuilt when
# Cargo.toml or Cargo.lock changes.
FROM rust:1.78-slim as planner
WORKDIR /app
RUN cargo install cargo-chef
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

# --- Stage 2: Builder ---
# This stage builds the actual application binary.
# It leverages the cached dependencies from the 'planner' stage.
FROM rust:1.78-slim as builder
WORKDIR /app
# First, copy the cached dependency recipe.
COPY --from=planner /app/recipe.json recipe.json
# Then, build the dependencies. This will be fast if the cache is valid.
RUN cargo chef cook --release --recipe-path recipe.json
# Now, copy the application source code and build it.
COPY . .
RUN cargo build --release --bin scoring-svc

# --- Stage 3: Final Image ---
# This stage creates the final, minimal image for deployment.
# We use a slim Debian base image because our application needs OpenSSL.
# For a truly minimal image, gcr.io/distroless/cc-debian12 is a great choice.
FROM debian:bookworm-slim as final

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app

# Copy the compiled binary from the 'builder' stage.
COPY --from=builder /app/target/release/scoring-svc /usr/local/bin/

# Set the user.
USER app

# Expose the gRPC port.
EXPOSE 50051

# The command to run when the container starts.
CMD ["/usr/local/bin/scoring-svc"]
