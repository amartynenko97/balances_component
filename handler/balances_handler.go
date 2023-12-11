package handler

import (
	"balances_component/constants"
	"balances_component/messaging"
	"balances_component/protofile"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

type BalancesHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
	logger            *logrus.Logger
	dbPool            *pgxpool.Pool
	messageQueue      chan Message
}

type Message struct {
	QueueType string
	Delivery  amqp.Delivery
}

func NewBalancesHandler(logger *logrus.Logger, publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel, connPool *pgxpool.Pool) *BalancesHandler {
	return &BalancesHandler{
		logger:            logger,
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
		dbPool:            connPool,
		messageQueue:      make(chan Message),
	}
}

func (h *BalancesHandler) StartListener(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		stopCh := make(chan struct{})
		defer close(stopCh)

		createAccountMessages := h.listeningChannel.ConsumeCreateAccountFromHttpApi(stopCh)
		getAccountBalancesMessages := h.listeningChannel.ConsumeGetAccountBalancesFromHttpApi(stopCh)

		for {
			select {
			case <-ctx.Done():
				return

			case delivery := <-createAccountMessages:
				h.messageQueue <- Message{QueueType: constants.QueueTypeCreateAccount, Delivery: delivery}

			case delivery := <-getAccountBalancesMessages:
				h.messageQueue <- Message{QueueType: constants.QueueTypeGetAccountBalances, Delivery: delivery}

			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:
		if err != nil {
			logrus.Println("Error in StartListener:", err)
		}
		return err
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
		return h.publishBalanceError(protofile.BalancesErrorCodes_BALANCE_ERROR_CODE_INTERNAL, "Internal error")
	}
	defer conn.Release()

	if err := createAccount(conn, &request); err != nil {
		logrus.WithError(err).Error("Error creating account")
		return h.publishBalanceError(protofile.BalancesErrorCodes_BALANCE_ERROR_CODE_INTERNAL, "Internal error")
	}

	return nil
}

func (h *BalancesHandler) publishBalanceError(errorCode protofile.BalancesErrorCodes, message string) error {
	errorMessage := &protofile.BalanceErrorMessage{}

	protoData, err := proto.Marshal(errorMessage)
	if err != nil {
		logrus.WithError(err).Error("Error marshalling balance error message")
		return err
	}

	err = h.publishingChannel.PublishCreateAccountToHttpApi(protoData)
	if err != nil {
		logrus.WithError(err).Error("Error publishing balance error message")
		return err
	}

	return nil
}
