FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./

# Build the Go application
RUN go build -o app .

# Expose the application port
EXPOSE 8100

# Command to run the built binary
CMD [ "./app" ]