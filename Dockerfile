# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o your-app ./cmd/server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/your-app .
COPY .env .

# Create data directory and set permissions
RUN mkdir /data && chmod -R 777 /data

VOLUME /data
EXPOSE 8080
CMD ["./your-app"]