package create_transaction

import (
	"context"
	"testing"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/entity"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/event"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/mocks"
	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("client1", "c1@c")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("client2", "c2@c")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	input := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()
	ctx := context.Background()

	c := NewCreateTransactionUseCase(
		mockUow,
		eventDispatcher,
		transactionCreatedEvent,
		balanceUpdatedEvent,
	)
	output, err := c.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
