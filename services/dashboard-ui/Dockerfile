# services/dashboard-ui/Dockerfile

# --- Stage 1: Dependency Installation ---
# This stage is dedicated to installing dependencies. Caching this layer
# speeds up builds if only the source code changes.
FROM node:18-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

# --- Stage 2: Builder ---
# This stage builds the Next.js application.
FROM node:18-alpine AS builder
WORKDIR /app
# Copy dependencies from the previous stage.
COPY --from=deps /app/node_modules ./node_modules
# Copy the application source code.
COPY . .
# Build the Next.js application for production.
RUN npm run build

# --- Stage 3: Final Image ---
# This stage creates the final, minimal image for production.
FROM node:18-alpine AS runner
WORKDIR /app

# Create a non-root user for security.
ENV NODE_ENV=production
RUN addgroup --system app && adduser --system --group app

# Copy the built application from the 'builder' stage.
# This includes the standalone output from Next.js.
COPY --from=builder /app/public ./public
COPY --from=builder --chown=app:app /app/.next/standalone ./
COPY --from=builder --chown=app:app /app/.next/static ./.next/static

# Switch to the non-root user.
USER app

# Expose the port the Next.js server runs on.
EXPOSE 3000

# The command to run when the container starts.
# This uses the Node.js server built by Next.js for production.
CMD ["node", "server.js"]
