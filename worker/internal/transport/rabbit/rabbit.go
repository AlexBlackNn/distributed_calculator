package rabbit

import (
	"context"
	"distributed_calculator/message_broker"
	"distributed_calculator/worker/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	connection        *amqp.Connection
	channel           *amqp.Channel
	operationQueue    amqp.Queue
	resultQueue       amqp.Queue
	logger            *slog.Logger
	calculatorService ServiceInterface
}

type ServiceInterface interface {
	Start(infix string) int
}

func New(log *slog.Logger, cfg *config.Config, calculatorService ServiceInterface) (*MessageBroker, error) {
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

	operationQueue, err := channel.QueueDeclare(
		"operation", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	resultQueue, err := channel.QueueDeclare(
		"result", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return nil, fmt.Errorf(
			"BROKER LAYER: broker.rabbit.New: couldn't create queue: %w",
			err,
		)
	}
	return &MessageBroker{
		connection:        connection,
		channel:           channel,
		operationQueue:    operationQueue,
		resultQueue:       resultQueue,
		logger:            log,
		calculatorService: calculatorService,
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

func (mb *MessageBroker) Send(ctx context.Context, message any, queueName string) error {
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
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
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

func (mb *MessageBroker) Receive(queueName string) error {

	messageChannel, err := mb.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"BROKER LAYER: broker.rabbit.Receive: couldn't get messageChannel: %w",
			err,
		)
	}

	var forever chan struct{}

	go func() {
		for msg := range messageChannel {
			ctx := context.Background()
			message := message_broker.Message{}
			err := json.Unmarshal(msg.Body, &message)
			fmt.Println(message)
			if err != nil {
				fmt.Println(err)
			}
			// TODO send result back using one more channel
			result := mb.calculatorService.Start(message.Operation)
			fmt.Println("=====>>>", mb.calculatorService.Start(message.Operation))
			mb.Send(ctx, result, "result")
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
