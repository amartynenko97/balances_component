package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountToHttpApi(protoData []byte) error
	PublishGetAccountBalancesToHttpApi(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountFromHttpApi(stopCh <-chan struct{}) <-chan amqp.Delivery
	ConsumeGetAccountBalancesFromHttpApi(stopCh <-chan struct{}) <-chan amqp.Delivery
}
