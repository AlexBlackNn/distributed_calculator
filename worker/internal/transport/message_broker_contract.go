package message_broker

import "github.com/google/uuid"

type ExectutionTime struct {
	PlusOperationExecutionTime           int
	MinusOperationExecutionTime          int
	MultiplicationOperationExecutionTime int
	Division_operation_execution_time    int
}

type RequestMessage struct {
	id                    uuid.UUID
	MessageExectutionTime ExectutionTime
	Operation             string
}

type ResultMessage struct {
	id    uuid.UUID
	Value int
	Err   error
}
