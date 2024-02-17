package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orchestrator/internal/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"orchestrator/message_broker"
)

type MessageBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func New(cfg *config.Config) (*MessageBroker, error) {
	connection, err := amqp.Dial(cfg.RabbitAmqp)
	if err != nil {
		return nil, fmt.Errorf(
			"BROKER LAYER: broker.rabbit.New: couldn't open a broker: %w",
			err,
		)
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf(
			"BROKER LAYER: broker.rabbit.New: couldn't create channel: %w",
			err,
		)
	}
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	queue, err := channel.QueueDeclare(
		"operation", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf(
			"BROKER LAYER: broker.rabbit.New: couldn't create queue: %w",
			err,
		)
	}
	return &MessageBroker{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}, nil
}

func (mb *MessageBroker) Stop() error {
	err := mb.connection.Close()
	if err != nil {
		return err
	}
	err = mb.channel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (mb *MessageBroker) Send(ctx context.Context, message message_broker.RequestMessage) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf(
			"BROKER LAYER: broker.rabbit.New: couldn't convert message %v to bytes: %w",
			message,
			err,
		)
	}
	err = mb.channel.PublishWithContext(ctx,
		"",            // exchange
		mb.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf(
			"BROKER LAYER: broker.rabbit.Send: couldn't send data: %w",
			err,
		)
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func (mb *MessageBroker) Receive() (chan message_broker.ResponseMessage, error) {

	results := make(chan message_broker.ResponseMessage)

	messageChannel, err := mb.channel.Consume(
		"result",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"BROKER LAYER: broker.rabbit.Receive: couldn't get messageChannel: %w",
			err,
		)
	}

	go func() {
		for msg := range messageChannel {
			message := message_broker.ResponseMessage{}
			err := json.Unmarshal(msg.Body, &message)
			fmt.Println(message)
			results <- message
			if err != nil {
				fmt.Println(err)
			}
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return results, nil
}
