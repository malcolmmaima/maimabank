# Use official Go builder image with Alpine Linux as the base
FROM golang:1.21rc3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Set GOPROXY to use the specified Go proxy
ENV GOPROXY=https://proxy.golang.org,direct

# Disable GOSUMDB to skip the checksum verification for modules
ENV GOSUMDB=off

# Fetch dependencies using Go Modules
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN --mount=type=cache,target=/go/pkg/mod go mod vendor

# Build the Go application with static linking
RUN --no-cache CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go
RUN --no-cache apk add curl
RUN --no-cache curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Use a minimal base image with Alpine Linux to reduce the image size
FROM alpine:3.13

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for .
COPY db/migration ./db/migration

# Install any necessary dependencies for your application
# For example, if your application requires SSL certificates, you may need to add them here.

# Expose port 8080 for the application
EXPOSE 8080

# Set the command to run the application
CMD ["./main"]
ENTRYPOINT ["/app/start.sh"]
