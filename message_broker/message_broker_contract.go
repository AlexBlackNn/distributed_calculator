package message_broker

type ExectutionTime struct {
	PlusOperationExecutionTime           int
	MinusOperationExecutionTime          int
	MultiplicationOperationExecutionTime int
	Division_operation_execution_time    int
}
type Message struct {
	MessageExectutionTime ExectutionTime
	Operation             string
}
