# Use the official Go image as the base image
FROM golang:1.22-alpine as build

# Set the working directory inside the container
WORKDIR /app

# Install the required dependencies
RUN apk add build-base

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Enable CGO to build the application
ENV CGO_ENABLED=1

# Build the Go application
RUN go build -o app ./cmd/main.go

# Use a minimal image to run the application
FROM alpine:3.13

# Copy the built binary from the previous stage
COPY --from=build /app/app .

ENTRYPOINT ["./app"]
