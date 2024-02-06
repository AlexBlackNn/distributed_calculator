package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"distributed_calculator/message_broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func New(brokerPath string) (*MessageBroker, error) {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

func (mb *MessageBroker) Send(ctx context.Context, message message_broker.Message) error {
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

func main() {
	ctx := context.Background()
	rabbitMqSender, err := New("test")
	if err != nil {
		fmt.Println(err)
	}

	execTime := message_broker.ExectutionTime{
		PlusOperationExecutionTime:           100,
		MinusOperationExecutionTime:          200,
		MultiplicationOperationExecutionTime: 300,
		Division_operation_execution_time:    400,
	}

	message := message_broker.Message{
		MessageExectutionTime: execTime,
		Operation:             "9*8+(7-8)",
	}

	err = rabbitMqSender.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
	}
}
