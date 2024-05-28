# Use the official Golang image to build the app
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o simplegowebserver main.go

# Start a new stage from scratch
FROM debian:stable-slim

# Set environment variables
ENV PORT=8080
ENV ADDR=0.0.0.0
ENV DIR=/app

# Add Maintainer info
LABEL maintainer="yourname@example.com"

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/simplegowebserver /usr/local/bin/simplegowebserver

# Copy the base directory content to the container
COPY --from=builder /app /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["simplegowebserver"]

