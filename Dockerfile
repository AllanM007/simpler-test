# Start from golang base image
FROM golang:1.23-alpine as builder

# Set necessary environmet variables needed
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the current working directory inside the container 
WORKDIR /app

# Copy the source from the current directory to the working directory inside the container 
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod tidy
RUN go mod download

# run tests
# RUN go test -v ./...

# build and run the project inside the container
RUN go build -o /app/cmd/main /app/cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the GIN_MODE environment variable
# ENV GIN_MODE=release

# Copy the .env file
COPY --from=builder /app/.env ./

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/main .

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["./main"]
