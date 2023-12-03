package balances

import (
	"balances_component/handler"
	"balances_component/logger"
	"balances_component/messaging"
	"golang.org/x/net/context"
	"log"
)

func main() {
	logger := logger.SetupLogger()

	rabbitMQConfig := messaging.RabbitMQConfig{
		URL: "amqp://guest:guest@localhost:5672/",
	}

	messageBroker, err := messaging.NewMessageBroker(rabbitMQConfig)
	if err != nil {
		log.Fatal("Failed to initialize MessageBroker:", err)
		return
	}

	defer messageBroker.Close()

	balancesHandler := handler.NewBalancesHandler(logger, messageBroker.GetPublishingChannel(), messageBroker.GetListeningChannel())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := balancesHandler.StartListener(ctx); err != nil {
			logger.Error("Error in StartListener:", err)

			cancel()
		}
	}()

}
