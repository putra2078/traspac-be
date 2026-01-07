package setup

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Declare(ch *amqp.Channel, exchange, exchangeType, queue, routingKey string) error {
	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	args := amqp.Table{
		"x-dead-letter-exchange": exchange + ".dlx",
	}

	if _, err := ch.QueueDeclare(
		queue,
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false,
		args,
	); err != nil {
		return err
	}

	return ch.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
}

// DeclareUserQueue - untuk per-user queue
func DeclareUserQueue(ch *amqp.Channel, exchange, exchangeType, userID string) error {
	queueName := "user." + userID + ".messages"
	routingKey := userID

	// Declare exchange (idempotent)
	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Declare user's queue
	args := amqp.Table{
		"x-message-ttl":          int32(86400000), // 24 hours
		"x-max-length":           int32(1000),     // max 1000 messages
		"x-dead-letter-exchange": exchange + ".dlx",
	}

	if _, err := ch.QueueDeclare(
		queueName,
		false, // durable
		true,  // auto-delete
		false, // exclusive (changed from true to allow multi-tab/reconnect)
		false, // no-wait
		args,
	); err != nil {
		return err
	}

	// Bind queue
	return ch.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
}
