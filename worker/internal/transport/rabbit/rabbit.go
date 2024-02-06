package rabbit

import (
	"context"
	"distributed_calculator/message_broker"
	"encoding/json"
	"fmt"
	"log"
	"time"

	transport "distributed_calculator/worker/internal/transport"
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

func (mb *MessageBroker) Send(ctx context.Context, message transport.Message) error {
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

func (mb *MessageBroker) Receive() error {

	messageChannel, err := mb.channel.Consume(
		mb.queue.Name,
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
			message := message_broker.Message{}
			err := json.Unmarshal(msg.Body, &message)
			fmt.Println(message)
			if err != nil {
				fmt.Println(err)
			}
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func main() {
	ctx := context.Background()
	rabbitMqSender, err := New("test")
	if err != nil {
		fmt.Println(err)
	}

	execTime := transport.ExectutionTime{
		PlusOperationExecutionTime:           100,
		MinusOperationExecutionTime:          200,
		MultiplicationOperationExecutionTime: 300,
		Division_operation_execution_time:    400,
	}
	message := transport.Message{
		MessageExectutionTime: execTime,
		Operation:             "9*8+(7-8)",
	}
	err = rabbitMqSender.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
	}
	rabbitMqSender.Receive()
}
