# Build stage
FROM golang:1.21.0 as builder

# Set the application directory
WORKDIR /app

# Enable Go modules
ENV GO111MODULE=on

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
COPY .env .
RUN go mod download

# Copy the application source
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o main .

# Execution stage
FROM gcr.io/distroless/base-debian10

# Copy the built binary
COPY --from=builder /app/main /

# Execute the application
CMD ["/main"]