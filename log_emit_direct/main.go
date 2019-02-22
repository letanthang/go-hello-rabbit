package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	if len(args) == 1 || args[1] == "" {
		return "hello"
	}

	return strings.Join(args[1:], " ")
}

func severityFrom(args []string) string {
	if len(args) == 1 || args[1] == "" {
		log.Fatalf("Usage: go run main.go [info|warning|error] hello_world")
	}

	s := args[1]
	log.Printf("Publish with routing key %s", s)

	return s
}

func main() {
	conn, err := amqp.Dial("amqp://rabbit:D66z3qm3ynC3@35.186.149.9")
	defer conn.Close()
	failOnError(err, "Fail to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", //name
		"direct",      //type
		true,          //durable
		false,         //auto delete
		false,         //internal
		false,         //nowait
		nil,
	)
	failOnError(err, "Failed to declare exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs_direct",         //exchange
		severityFrom(os.Args), //routing key
		false, //mandatory
		false, //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Send a msg: %s", body)
}
