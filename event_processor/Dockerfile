# Build stage
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod init event_processor
RUN go get github.com/go-redis/redis/v8 github.com/streadway/amqp
RUN go build -o event_processor .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/event_processor .
CMD ["./event_processor"]
