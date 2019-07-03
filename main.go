package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	rabbitHost := "amqp://guest:guest@localhost:5672"

	conn, err := amqp.Dial(rabbitHost)
	failOnError(err, "Failed to connect RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	args := make(amqp.Table)
	args["x-max-priority"] = int64(10)

	q, err := ch.QueueDeclare(
		"tix.order.queue",
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	ch.QueueBind(
		q.Name,
		"",
		"amq.fanout",
		false,
		nil,
	)

	msg := amqp.Publishing{
		Body: []byte("my first message"),
	}

	ch.Publish(
		"",
		q.Name,
		false,
		false,
		msg,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	for m := range msgs {
		println(string(m.Body))
	}
}
