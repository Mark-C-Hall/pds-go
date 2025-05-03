# Start from the official Go image
FROM golang:1.24-alpine AS build

# Set working directory inside the container
WORKDIR /app

# Copy the entire source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/pds/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Copy the binary from build stage
COPY --from=build /app/main /app/main

# Set working directory
WORKDIR /app

# Command to run the executable
CMD ["./main"]