package connection

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(url string) (*amqp.Connection, error) {
	return amqp.DialConfig(url, amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "id_ID",
	})
}
