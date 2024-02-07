package message_broker

import (
	"context"
)

type MessageBrokerInterface interface {
	Receive() error
	Send(context.Context, RequestMessage) error
	Stop() error
}
