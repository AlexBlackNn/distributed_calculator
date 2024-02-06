package main

import (
	"context"
	transport "distributed_calculator/worker/internal/transport"
	rabbit "distributed_calculator/worker/internal/transport/rabbit"
	"fmt"
)

func main() {
	ctx := context.Background()
	rabbitMqSender, err := rabbit.New("test")
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
