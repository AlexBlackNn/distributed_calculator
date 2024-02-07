package rabbit

import (
	"context"
	"distributed_calculator/worker/internal/config"
	transport "distributed_calculator/worker/internal/transport"
	"encoding/json"
	"fmt"
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
	connection, err := amqp.Dial(cfg.Rabbit.Amqp)
	if err != nil {
		return nil, fmt.Errorf(
			"TRANSPORT LAYER: rabbit.New: couldn't open a broker: %w",
			err,
		)
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf(
			"TRANSPORT LAYER: rabbit.New: couldn't create channel: %w",
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
			"TRANSPORT LAYER: rabbit.New: couldn't create queue: %w",
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

func (mb *MessageBroker) Send(ctx context.Context, message any, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log := mb.logger.With(
		slog.String("info", "TRANSPORT LAYER: Send"),
	)
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf(
			"TRANSPORT LAYER: rabbit.New: couldn't convert message %v to bytes: %w",
			message,
			err,
		)
	}
	log.Info("Marshal message %v", message)

	err = mb.channel.PublishWithContext(ctx,
		"",                    // exchange
		cfg.Rabbit.WriteQueue, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf(
			"TRANSPORT LAYER: rabbit.Send: couldn't send data: %w",
			err,
		)
	}
	log.Info("Publish message %v", message)
	return nil
}

func (mb *MessageBroker) Receive(cfg *config.Config) error {
	log := mb.logger.With(
		slog.String("info", "TRANSPORT LAYER: Receive"),
	)

	messageChannel, err := mb.channel.Consume(
		cfg.Rabbit.ReadQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"TRANSPORT LAYER: rabbit.Receive: couldn't get messageChannel: %w",
			err,
		)
	}
	log.Info("Receiver is ready!")
	var forever chan struct{}

	go func() {
		for msg := range messageChannel {
			requestMessage := transport.RequestMessage{}
			err := json.Unmarshal(msg.Body, &requestMessage)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(requestMessage)
			result := mb.calculatorService.Start(requestMessage.Operation)
			responseMessage := transport.ResponseMessage{
				Id:    requestMessage.Id,
				Value: result,
				Err:   nil,
			}
			fmt.Println("=====>>>", responseMessage)
			mb.Send(context.Background(), responseMessage, cfg)
			msg.Ack(false)
		}
	}()
	log.Info("[*] Waiting for messages!")
	<-forever
	return nil
}
