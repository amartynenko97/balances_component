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

func (l *Listener) ConsumeCreateAccountFromHttpApi(stopCh <-chan struct{}) <-chan amqp.Delivery {
	return l.consume(constants.CreateAccountRequestQueue, stopCh)
}

func (l *Listener) ConsumeGetAccountBalancesFromHttpApi(stopCh <-chan struct{}) <-chan amqp.Delivery {
	return l.consume(constants.GetAccountBalanceRequestQueue, stopCh)
}

func (l *Listener) consume(queueName string, stopCh <-chan struct{}) <-chan amqp.Delivery {
	messages := make(chan amqp.Delivery)

	go func() {
		defer close(messages)

		consumer, err := l.channel.Consume(
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

		for {
			select {
			case <-stopCh:
				return
			case delivery, ok := <-consumer:
				if !ok {
					return
				}

				select {
				case messages <- delivery:
				case <-stopCh:
					return
				}
			}
		}
	}()

	return messages
}
