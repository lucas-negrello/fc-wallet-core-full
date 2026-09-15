package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/create_transaction"
)

type TransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (handler *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := handler.CreateTransactionUseCase.Execute(r.Context(), dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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
