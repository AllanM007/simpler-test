# Start from golang base image
FROM golang:1.23-alpine as builder

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download && go mod verify

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64 

# Change to the cmd directory where the main.go is located
WORKDIR /app/cmd

RUN go build -o /go/bin/simpler-test/cmd/main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Set the GIN_MODE environment variable
ENV GIN_MODE=release

# Copy the .env file (if necessary for runtime)
COPY --from=builder /app/.env ./

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /go/bin/simpler-test/cmd/main .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]
