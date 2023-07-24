# Use golang:1.21rc3 as the base image for building the Go application
FROM golang:1.21rc3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Set GOPROXY to use the specified Go proxy
ENV GOPROXY=https://proxy.golang.org,direct

# Disable GOSUMDB to skip the checksum verification for modules
ENV GOSUMDB=off

# Fetch dependencies using Go Modules
RUN go mod download
RUN go mod vendor

# Build the Go application
RUN go build -o main main.go

# Use a minimal base image to reduce the image size
FROM scratch

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/main .

# Expose port 8080 for the application
EXPOSE 8080

# Set the command to run the application
CMD ["./main"]
