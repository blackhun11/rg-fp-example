FROM golang:latest

# Set the working directory to the project root
WORKDIR /go/src/app

# Copy the project source code
COPY . .

# Download and install dependencies
RUN go get -d -v ./...

# Build the application
RUN go build

# Set the entrypoint to the application binary
CMD ["./wb"]

