package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/create_account"
)

type AccountHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
}

func NewAccountHandler(createAccountUseCase create_account.CreateAccountUseCase) *AccountHandler {
	return &AccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (handler *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto create_account.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := handler.CreateAccountUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
