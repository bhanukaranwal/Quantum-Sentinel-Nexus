# services/auth-gateway-svc/Dockerfile

# --- Stage 1: Builder ---
# This stage installs all dependencies, including devDependencies,
# and builds the TypeScript source code into JavaScript.
FROM node:18-alpine AS builder

# Set the working directory.
WORKDIR /app

# Copy the package.json and package-lock.json files.
COPY package*.json ./

# Install all dependencies.
RUN npm install

# Copy the rest of the application source code.
COPY . .

# Build the application. This runs the "build" script from package.json,
# which compiles TypeScript to JavaScript in the 'dist' directory.
RUN npm run build

# Remove development dependencies to lighten the next stage.
RUN npm prune --production


# --- Stage 2: Final Image ---
# This stage creates the final, minimal image for production.
FROM node:18-alpine

# Set the working directory.
WORKDIR /app

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app

# Copy the production node_modules and the compiled 'dist' directory
# from the 'builder' stage.
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./

# Switch to the non-root user.
USER app

# Expose the port the service listens on.
EXPOSE 3001

# The command to run when the container starts.
# This runs the compiled JavaScript from the 'dist' directory.
CMD [ "node", "dist/main" ]
