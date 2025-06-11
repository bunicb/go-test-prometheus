# Stage 1: Build the Go application
# We use a specific Go version for reproducibility.
# Alpine Linux is used for a smaller image size.
FROM docker.io/golang:1.22-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
# This leverages Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY *.go ./

# Build the application. CGO_ENABLED=0 is important for a static binary
# that can run in a minimal container without C libraries.
# GOOS=linux specifies the target operating system.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /main .

# Stage 2: Create the final, minimal image
# We use a minimal base image, `alpine`, for the final container.
FROM docker.io/alpine:latest

# It's good practice to run containers as a non-root user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy only the compiled binary from the builder stage.
COPY --from=builder /main /main

# Expose the port the application runs on.
EXPOSE 8080

# The command to run when the container starts.
CMD ["/main"]

