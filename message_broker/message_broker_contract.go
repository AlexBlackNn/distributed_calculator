package message_broker

type ExectutionTime struct {
	PlusOperationExecutionTime           int
	MinusOperationExecutionTime          int
	MultiplicationOperationExecutionTime int
	DivisionOperationExecutionTime       int
}

type RequestMessage struct {
	Id                    string
	MessageExectutionTime ExectutionTime
	Operation             string
}

type ResponseMessage struct {
	Id    string
	Value int
	Err   error
}
