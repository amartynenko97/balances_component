package handler

import (
	"balances_component/messaging"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"log"
)

type BalancesHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
	logger            *logrus.Logger
}

func NewBalancesHandler(logger *logrus.Logger, publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel) *BalancesHandler {
	return &BalancesHandler{
		logger:            logger,
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
	}
}

func (h *BalancesHandler) StartListener(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- h.StartConsumer(ctx, errCh)
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:
		if err != nil {
			log.Println("Error in StartConsumer:", err)
		}
		return err
	}
}

func (h *BalancesHandler) StartConsumer(ctx context.Context, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery := <-h.listeningChannel.ConsumeCreateAccountBalancesFromHttpApi():
			go func(delivery amqp.Delivery) {
				if err := h.processCreateAccountBalance(delivery.Body); err != nil {
					logrus.Infoln("Error in processCreateAccountBalance:", err)
					errCh <- err
				}
			}(delivery)

		case delivery := <-h.listeningChannel.ConsumeGetAccountBalancesFromHttpApi():
			go func(delivery amqp.Delivery) {
				if err := h.processGetAccountBalance(delivery.Body); err != nil {
					logrus.Infoln("Error in processGetAccountBalance:", err)
					errCh <- err
				}
			}(delivery)
		}
	}
}

func (h *BalancesHandler) processCreateAccountBalance(protoData []byte) error {
	// Ваш код обработки CreateAccountBalance
	// ...
	// Если произошла ошибка
	return errors.New("some error occurred")
}

func (h *BalancesHandler) processGetAccountBalance(protoData []byte) error {
	// Ваш код обработки GetAccountBalance
	// ...
	// Если произошла ошибка
	return errors.New("some error occurred")
}

//func (h *BalancesHandler) processCreateAccountBalance(protoData []byte) {
//	errorResponse := constants.ErrorResponse{
//		Error: string(constants.NoSuchCurrency),
//	}
//
//	errorJSON, err := json.Marshal(errorResponse)
//	if err != nil {
//		logrus.Errorf("Error marshalling error response: %v", err)
//		return
//	}
//
//	err = h.publishingChannel.PublishCreateAccountBalancesToHttpApi(errorJSON)
//	if err != nil {
//		logrus.Errorf("Error publishing error response: %v", err)
//		return
//	}
//}
