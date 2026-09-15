package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/database"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/event"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/event/handler"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/create_account"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/create_client"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com.br/lucas-negrello/fc-ms-wallet/internal/web"
	webserver2 "github.com.br/lucas-negrello/fc-ms-wallet/internal/web/webserver"
	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/events"
	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/kafka"
	"github.com.br/lucas-negrello/fc-ms-wallet/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)
import _ "github.com/go-sql-driver/mysql"

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/wallet?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))

	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()
	//eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)
	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver2.NewWebServer(":8080")
	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Print("Server running on port 8080")
	webserver.Start()
}
