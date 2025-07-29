# services/compliance-svc/Dockerfile

# --- Stage 1: Builder ---
# This stage uses an official Eclipse Temurin image with the full JDK,
# which is required to build the Gradle project.
FROM eclipse-temurin:17-jdk-jammy as builder

# Set the working directory.
WORKDIR /app

# Copy the Gradle wrapper and build script.
COPY gradlew .
COPY gradle gradle
COPY build.gradle.kts .
COPY settings.gradle.kts .

# Run the bootJar task to build the executable JAR.
# The --no-daemon flag is recommended for CI/CD environments.
RUN ./gradlew bootJar --no-daemon

# --- Stage 2: Final Image ---
# This stage creates the final, minimal image for deployment.
# We use a JRE image, which is smaller than a full JDK image.
FROM eclipse-temurin:17-jre-jammy

# Set the working directory.
WORKDIR /app

# Create a non-root user for security.
RUN addgroup --system app && adduser --system --group app
USER app

# Copy the executable JAR file from the 'builder' stage.
# The JAR file is located in the build/libs directory.
COPY --from=builder /app/build/libs/*.jar app.jar

# Expose the port the Spring Boot application listens on.
EXPOSE 8083

# The command to run when the container starts.
# This executes the Spring Boot application.
ENTRYPOINT ["java", "-jar", "app.jar"]
