package message_broker

type ExectutionTime struct {
	PlusOperationExecutionTime           int
	MinusOperationExecutionTime          int
	MultiplicationOperationExecutionTime int
	DivisionOperationExecutionTime       int
}

type RequestMessage struct {
	Id                   string
	MessageExecutionTime ExectutionTime
	Operation            string
}

type ResultMessage struct {
	Id    string
	Value int
	Err   error
}
