package producer

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ch *amqp.Channel, exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}


// PublishToUser - untuk publish ke specific user
func PublishToUser(ch *amqp.Channel, exchange, recipientID string, body []byte, senderID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return ch.PublishWithContext(
        ctx,
        exchange,
        recipientID, // routing key = recipient's userID
        false,
        false,
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/json",
            Body:         body,
            Timestamp:    time.Now(),
            Headers: amqp.Table{
                "sender": senderID,
            },
        },
    )
}