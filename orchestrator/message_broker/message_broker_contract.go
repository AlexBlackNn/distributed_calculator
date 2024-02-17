package message_broker

import "orchestrator/internal/domain/models"

type RequestMessage struct {
	Id                    string
	MessageExectutionTime models.Settings
	Operation             string
}

type ResponseMessage struct {
	Id    string
	Value int
	Err   string
}
