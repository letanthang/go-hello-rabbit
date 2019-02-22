package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://rabbit:D66z3qm3ynC3@35.186.149.9")
	failOnError(err, "Failed to connect RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Fail to create a channel")

	q, err := ch.QueueDeclare(
		"",    //queue name
		true,  //durable
		false, //delete when used
		true,  //exclusive
		false, //no wait
		nil,
	)
	failOnError(err, "Fail to declare queue")

	ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)

	failOnError(err, "Fail to bind a queue")

	err = ch.Qos(
		1,     //prefetch count
		0,     //prefetch size
		false, //global
	)
	failOnError(err, "Failed to set qos")

	msgs, err := ch.Consume(
		q.Name, //queue name
		"",     //consumer
		true,   //auto ack
		false,  //exclusive
		false,  //no local
		false,  //no wait
		nil,
	)

	failOnError(err, "Fail to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// dot_count := bytes.Count(d.Body, []byte("."))
			// log.Printf("Dot count sleep time: %s", dot_count)
			// t := time.Duration(dot_count)
			// time.Sleep(t * time.Second)
			// log.Printf("Done")
			// d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit presss ctrl+c")
	<-forever
}
