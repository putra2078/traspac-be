package worker

import (
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Start(
	wg *sync.WaitGroup,
	jobs <-chan amqp.Delivery,
	id int,
) {
	defer wg.Done()

	for msg := range jobs {
		log.Printf("[worker %d] processing message", id)

		if err := process(msg.Body); err != nil {
			log.Printf("[worker %d] error: %v", id, err)
			msg.Nack(false, false) // no requeue → DLQ
			continue
		}

		msg.Ack(false)
	}
}

func process(data []byte) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}
