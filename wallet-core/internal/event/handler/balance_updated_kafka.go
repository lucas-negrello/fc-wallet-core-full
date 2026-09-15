package handler

import (
	"fmt"
	"sync"

	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/events"
	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/kafka"
)

type UpdateBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *UpdateBalanceKafkaHandler {
	return &UpdateBalanceKafkaHandler{
		Kafka: kafka,
	}
}

func (h *UpdateBalanceKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	err := h.Kafka.Publish(message, nil, "balances")
	if err != nil {
		return
	}
	fmt.Println("BalanceUpdatedKafkaHandler: ", message.GetPayload())
}
