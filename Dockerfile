FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /app/config/config.yaml .

RUN apk --no-cache add tzdata
RUN apk --no-cache add curl
ENV TZ=Asia/Bangkok
ENV API_CONFIG_PATH=/root
ENV API_CONFIG_NAME=config

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
