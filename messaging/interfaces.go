package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountBalancesToHttpApi(protoData []byte) error
	PublishGetAccountBalancesToHttpApi(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountBalancesFromHttpApi() <-chan amqp.Delivery
	ConsumeGetAccountBalancesFromHttpApi() <-chan amqp.Delivery
}
