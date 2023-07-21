# Use the official Alpine Linux image as the base
FROM alpine:3.13

# Install necessary packages (including Git) and set the Alpine package mirror
RUN apk update && apk add --no-cache git && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.13/main" > /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.13/community" >> /etc/apk/repositories

# Set custom DNS servers for the container
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
    echo "nameserver 8.8.4.4" >> /etc/resolv.conf

# Use golang:1.21rc3 as the base image for the Go application
FROM golang:1.21rc3

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Set GOPROXY to use proxy.golang.org as the Go proxy
ENV GOPROXY=https://proxy.golang.org,direct

# Disable GOSUMDB to skip the checksum verification for modules
ENV GOSUMDB=off

# Fetch dependencies using Go Modules
RUN go mod download
RUN go mod vendor

# Build the Go application
RUN go build -o main main.go

# Expose port 8080 for the application
EXPOSE 8080

# Set the command to run the application
CMD ["./main"]
