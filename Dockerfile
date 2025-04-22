FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ham-radio-assistant ./cmd/ham-radio-assistant

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the compiled executable
COPY --from=builder /app/ham-radio-assistant /app/
# Copy config file
COPY config.json /app/

# Expose the port defined in config.json
EXPOSE 8080

# Run the executable
CMD ["/app/ham-radio-assistant"]
