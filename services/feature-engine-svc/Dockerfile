# services/feature-engine-svc/Dockerfile

# --- Build Stage ---
# This stage installs the Python dependencies. Using a separate stage
# allows Docker to cache this layer, so dependencies are only re-installed
# when requirements.txt changes, speeding up subsequent builds.
FROM python:3.12-slim as builder

# Set the working directory
WORKDIR /app

# Install build dependencies that might be needed for some Python packages.
RUN apt-get update && apt-get install -y --no-install-recommends build-essential

# Copy the requirements file and install the dependencies.
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt


# --- Final Stage ---
# This stage creates the final application image.
# We use a slim base image to keep the size down.
FROM python:3.12-slim

# Set the working directory
WORKDIR /app

# Create a non-root user to run the application for security.
RUN addgroup --system app && adduser --system --group app

# Copy the installed dependencies from the 'builder' stage.
COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# Copy the application source code.
COPY . .

# Change ownership of the app directory to the non-root user.
RUN chown -R app:app /app

# Switch to the non-root user.
USER app

# The command to run when the container starts.
CMD ["python", "main.py"]
