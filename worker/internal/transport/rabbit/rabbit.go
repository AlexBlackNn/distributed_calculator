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
	Start(transport.RequestMessage) int
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

func (mb *MessageBroker) Stop() {
	log := mb.logger.With(
		slog.String("info", "TRANSPORT LAYER: Stop"),
	)
	err := mb.connection.Close()
	if err != nil {
		log.Error(
			"TRANSPORT LAYER: rabbit.Stop: couldn't close rabbit connection",
			"error", err.Error(),
		)
	}
	err = mb.channel.Close()
	if err != nil {
		log.Error(
			"TRANSPORT LAYER: rabbit.Stop: couldn't close rabbit channel",
			"error", err.Error(),
		)
	}
}

func (mb *MessageBroker) Send(ctx context.Context, message any, cfg *config.Config) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log := mb.logger.With(
		slog.String("info", "TRANSPORT LAYER: Send"),
	)
	body, err := json.Marshal(message)
	if err != nil {
		log.Error(
			"TRANSPORT LAYER: rabbit.Send: couldn't convert message to bytes",
			"message", message, "error", err.Error(),
		)
	}
	log.Info("Marshal message: ", "message", message)

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
		log.Error(
			"TRANSPORT LAYER: rabbit.Send: couldn't send message",
			"message", message, "error", err.Error(),
		)
	}
	log.Info("Publish message: ", "message", message)
}

func (mb *MessageBroker) Receive(cfg *config.Config) {
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
		log.Error(
			"TRANSPORT LAYER: rabbit.Receive: couldn't get messageChannel",
			"error", err.Error(),
		)
	}
	log.Info("Receiver is ready!")
	var forever chan struct{}

	go func() {
		for msg := range messageChannel {
			requestMessage := transport.RequestMessage{}
			err := json.Unmarshal(msg.Body, &requestMessage)
			if err != nil {
				log.Error(err.Error())
			}
			log.Info("get request message: ", "message", requestMessage)
			result := mb.calculatorService.Start(requestMessage)
			responseMessage := transport.ResponseMessage{
				Id:    requestMessage.Id,
				Value: result,
				Err:   nil,
			}
			log.Info("formed response message", "message", responseMessage)
			mb.Send(context.Background(), responseMessage, cfg)
			msg.Ack(false)
		}
	}()
	log.Info("[*] Waiting for messages!")
	<-forever
}
