package messaging

import (
	"balances_component/constants"
	"github.com/streadway/amqp"
	"log"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{
		channel: channel,
	}
}

func (p *Publisher) PublishCreateAccountToHttpApi(protoData []byte) error {
	return p.publish(protoData, constants.ExNameBalances, constants.RkCreateAccountResponse)
}

func (p *Publisher) PublishGetAccountBalancesToHttpApi(protoData []byte) error {
	return p.publish(protoData, constants.ExNameBalances, constants.RkGetAccountBalanceResponse)
}

func (p *Publisher) publish(protoData []byte, exchange string, routingKey string) error {
	err := p.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        protoData,
		},
	)

	if err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Printf("Message published successfully with routing key: %s\n", routingKey)
	return nil
}
