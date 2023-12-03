package messaging

import (
	"balances_component/constants"
	"github.com/streadway/amqp"
	"log"
)

type Listener struct {
	channel *amqp.Channel
}

func NewListener(channel *amqp.Channel) *Listener {
	return &Listener{
		channel: channel,
	}
}

func (l *Listener) ConsumeCreateAccountBalancesFromHttpApi() <-chan amqp.Delivery {
	return l.consume(constants.CreateAccountBalanceRequestQueue)
}

func (l *Listener) ConsumeGetAccountBalancesFromHttpApi() <-chan amqp.Delivery {
	return l.consume(constants.GetAccountBalanceRequestQueue)
}

func (l *Listener) consume(queueName string) <-chan amqp.Delivery {
	messages, err := l.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	return messages
}
