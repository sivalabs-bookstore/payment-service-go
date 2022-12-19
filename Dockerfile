FROM golang:1.19-buster as builder
# Create and change to the app directory.
WORKDIR /app
# Copy go.mod and if present go.sum.
COPY go.* ./
# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy local code to the container image.
COPY . ./
# Build the Go app
RUN GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -v -o payment-service


######## Start a new stage from scratch #######
FROM gcr.io/distroless/base-debian10
WORKDIR /

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/payment-service ./payment-service
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Run the templates service on container startup.
CMD ["/payment-service"]
