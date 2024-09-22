package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type EntryEvent struct {
	ID            string `json:"id"`
	VehiclePlate  string `json:"vehicle_plate"`
	EntryDateTime string `json:"entry_date_time"`
}

type ExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

var ctx = context.Background()

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	entryQ, err := ch.QueueDeclare(
		"vehicle_entry_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare entry queue")

	exitQ, err := ch.QueueDeclare(
		"vehicle_exit_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare exit queue")

	msgsEntry, err := ch.Consume(
		entryQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register entry consumer")

	msgsExit, err := ch.Consume(
		exitQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register exit consumer")

	go func() {
		for d := range msgsEntry {
			var entry EntryEvent
			json.Unmarshal(d.Body, &entry)

			err := redisClient.Set(ctx, entry.VehiclePlate+":entry", entry.EntryDateTime, 0).Err()
			if err != nil {
				log.Printf("Error setting entry in Redis: %s", err)
			}

			log.Printf("Processed entry event: %s", entry.VehiclePlate)
		}
	}()

	go func() {
		for d := range msgsExit {
			var exit ExitEvent
			json.Unmarshal(d.Body, &exit)

			entryTime, err := redisClient.Get(ctx, exit.VehiclePlate+":entry").Result()
			if err == redis.Nil {
				log.Printf("No entry found for exit event: %s", exit.VehiclePlate)
				continue
			} else if err != nil {
				log.Printf("Error getting entry from Redis: %s", err)
				continue
			}

			entryTimeParsed, _ := time.Parse(time.RFC3339, entryTime)
			exitTimeParsed, _ := time.Parse(time.RFC3339, exit.ExitDateTime)
			duration := exitTimeParsed.Sub(entryTimeParsed)

			log.Printf("Vehicle: %s, Entry: %s, Exit: %s, Duration: %s",
				exit.VehiclePlate, entryTime, exit.ExitDateTime, duration)

			// TODO: Send parking summary to the REST API

			redisClient.Del(ctx, exit.VehiclePlate+":entry")
		}
	}()

	log.Println("Event Processor Running")
	select {}
}
