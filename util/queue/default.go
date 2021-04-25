package queue

import (
	"fmt"
	"github.com/streadway/amqp"
)

type DefaultQueueService struct {
	AMQPServerURL string
}

func NewDefaultQueueService(amqpServerURL string) (defaultQueueService Service) {
	defaultQueueService = &DefaultQueueService{
		AMQPServerURL: amqpServerURL,
	}
	return
}

func (q *DefaultQueueService) PublishOrderQueue(productID uint64, quantity uint32) {
	// Define RabbitMQ server URL.
	amqpServerURL := q.AMQPServerURL

	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer conn.Close()

	// Opening a channel to our RabbitMQ instance over
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer ch.Close()

	// Publish
	queue, err := ch.QueueDeclare(
		"order_product",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(queue)

	// Publish
	msg := fmt.Sprintf("%d:%d", productID, quantity)
	err = ch.Publish(
		"",
		"order_product",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Successfully published message to Queue")
}
