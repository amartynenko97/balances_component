package balances

import (
	"balances_component/database"
	"balances_component/handler"
	"balances_component/logger"
	"balances_component/messaging"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logger.SetupLogger()

	rabbitMQConfig := messaging.RabbitMQConfig{
		URL: "amqp://admin:admin@rabbitmq:5672",
	}

	messageBroker, err := messaging.NewMessageBroker(rabbitMQConfig)
	if err != nil {
		logger.Fatal("Failed to initialize MessageBroker:", err)
		return
	}
	defer messageBroker.Close()

	dbPool, err := database.GetConnection(context.Background(), "postgres://postgre:admin@postgres:5432/balances", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize PostgreSQL connection pool")
		return
	}
	defer dbPool.Close()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	balancesHandler := handler.NewBalancesHandler(logger, messageBroker.GetPublishingChannel(), messageBroker.GetListeningChannel(), dbPool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := balancesHandler.StartListener(ctx); err != nil {
			logger.WithError(err).Fatal("Error in StartListener:", err)

			cancel()
		}
	}()

	<-shutdownChan
	cancel()
	logger.Info("Shutting down gracefully...")
}
