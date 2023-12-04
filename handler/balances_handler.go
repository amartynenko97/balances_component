package handler

import (
	"balances_component/messaging"
	"balances_component/protofile"
	"errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type BalancesHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
	logger            *logrus.Logger
	dbPool            *pgxpool.Pool
}

func NewBalancesHandler(logger *logrus.Logger, publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel, connPool *pgxpool.Pool) *BalancesHandler {
	return &BalancesHandler{
		logger:            logger,
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
		dbPool:            connPool,
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
			logrus.Println("Error in StartConsumer:", err)
		}
		return err
	}
}

func (h *BalancesHandler) StartConsumer(ctx context.Context, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery := <-h.listeningChannel.ConsumeCreateAccountFromHttpApi():
			go func(delivery amqp.Delivery) {
				if err := h.processCreateAccount(delivery.Body); err != nil {
					errCh <- err
				}
			}(delivery)

		case delivery := <-h.listeningChannel.ConsumeGetAccountBalancesFromHttpApi():
			go func(delivery amqp.Delivery) {
				if err := h.processGetAccountBalance(delivery.Body); err != nil {
					errCh <- err
				}
			}(delivery)
		}
	}
}

func (h *BalancesHandler) processCreateAccount(protoData []byte) error {
	var request protofile.CreateAccountRequest

	if err := jsonpb.UnmarshalString(string(protoData), &request); err != nil {
		logrus.WithError(err).Error("Error unmarshalling proto data")
		return err
	}
	request.UserId = uuid.New().String()

	conn, err := h.dbPool.Acquire(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Error acquiring database connection")
		return err
	}
	defer conn.Release()

	if err := createAccount(conn, &request); err != nil {
		logrus.WithError(err).Error("Error creating account")
		return err
	}

	return nil
}

func (h *BalancesHandler) processGetAccountBalance(protoData []byte) error {
	return errors.New("some error occurred")
}
