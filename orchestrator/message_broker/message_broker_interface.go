package message_broker

import (
	"context"
)

type MessageBrokerInterface interface {
	Receive() (chan ResponseMessage, error)
	Send(context.Context, RequestMessage) error
	Stop() error
}
