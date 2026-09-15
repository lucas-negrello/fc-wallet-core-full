package create_account

import (
	"testing"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/entity"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j")
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)

	c := NewCreateAccountUseCase(accountMock, clientMock)
	input := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	output, err := c.Execute(input)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.ID)
	clientMock.AssertExpectations(t)
	accountMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
