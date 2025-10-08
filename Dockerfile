# Dockerfile for Windows Go cross-compilation on Mac
FROM golang:1.23-bullseye

# Install necessary packages for Windows cross-compilation
RUN apt-get update && apt-get install -y \
    gcc-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Set environment variables for Windows cross-compilation
ENV GOOS=windows
ENV GOARCH=amd64
ENV CGO_ENABLED=1
ENV CC=x86_64-w64-mingw32-gcc

# Copy go.mod and go.sum (if exists)
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Build the Windows executable
RUN go build -o WingStealer.exe .

# Keep container running for debugging
CMD ["tail", "-f", "/dev/null"]
