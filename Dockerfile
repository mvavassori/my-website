# Start from the official Golang image.
FROM golang:alpine

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache.
COPY go.mod go.sum ./

# Download necessary Go modules.
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the application.
RUN go build -o my-website .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Run the binary.
CMD ["./my-website"]