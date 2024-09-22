package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

type ExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"vehicle_exit_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for {
		event := ExitEvent{
			ID:           fmt.Sprintf("%d", rand.Intn(100000)),
			VehiclePlate: fmt.Sprintf("VEH-%d", rand.Intn(9999)),
			ExitDateTime: time.Now().UTC().Format(time.RFC3339),
		}

		body, err := json.Marshal(event)
		failOnError(err, "Failed to marshal JSON")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		failOnError(err, "Failed to publish a message")

		log.Printf("Sent exit event: %s", body)
		time.Sleep(2 * time.Second) // Simulate exit every 2 seconds
	}
}
