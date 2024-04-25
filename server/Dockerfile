# Start from the latest golang base image
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . .

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# Build the Go app
RUN go build -ldflags="-s -w" -o app main.go


# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/app .

RUN chmod +x app

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./app"]
