# # use official Golang image
# FROM golang:1.16.3-alpine3.13

# # set working directory
# WORKDIR /app

# # Copy the source code
# COPY . . 

# # Download and install the dependencies
# RUN go get -d -v ./...

# # Build the Go app
# RUN go build -o api .

# #EXPOSE the port
# EXPOSE 8000

# # Run the executable
# CMD ["./api"]


# Build stage
FROM golang:1.22-alpine AS builder

# Enable Go modules and configure working dir
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /app

# Copy go.mod and go.sum first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o api .

# Final stage: minimal image
FROM alpine:latest

# Set up non-root user (optional, security best practice)
RUN adduser -D appuser
USER appuser

WORKDIR /app
COPY --from=builder /app/api .

EXPOSE 8000

CMD ["./api"]
