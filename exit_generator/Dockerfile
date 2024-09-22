# Build stage
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod init exit_generator
RUN go build -o exit_generator .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/exit_generator .
CMD ["./exit_generator"]
