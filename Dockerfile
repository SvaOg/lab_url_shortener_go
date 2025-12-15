# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies (Docker will cache it)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
# CGO_ENABLED=0 creates statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Start (minimal image)
FROM alpine:latest

WORKDIR /root/

# Copy binaries from the first stage
COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]