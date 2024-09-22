# Build stage
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod init entry_generator
RUN go build -o entry_generator .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/entry_generator .
CMD ["./entry_generator"]
