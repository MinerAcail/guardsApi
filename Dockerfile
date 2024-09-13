# Use the official Go image as the base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Install Air
# RUN go install github.com/cosmtrek/air@latest no more used in the docs
RUN go install github.com/air-verse/air@latest


# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Create tmp directory
RUN mkdir -p /app/tmp

# Copy the Air configuration file
COPY .air.toml ./

# Build the Go app
RUN go build -o /app/tmp/main -buildvcs=false .

# Command to run Air
CMD ["air"]
