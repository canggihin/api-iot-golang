# Stage 1: Build the application
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .
RUN go build -v -o main .

# Stage 2: Prepare CA certificates using Alpine
FROM alpine AS certs
RUN apk --no-cache add ca-certificates

# Stage 3: Create the final Docker image with scratch
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy CA certificates from the certs stage
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8089

CMD ["./main"]
