package consumer

import (
	"context"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"hrm-app/internal/pkg/rabbitmq/config"
	"hrm-app/internal/pkg/rabbitmq/worker"
)

func Start(ctx context.Context, ch *amqp.Channel) error {
	if err := ch.Qos(config.Prefetch, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		config.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	jobs := make(chan amqp.Delivery, config.WorkerCount)
	wg := sync.WaitGroup{}

	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go worker.Start(&wg, jobs, i)
	}

	for {
		select {
		case msg := <-msgs:
			jobs <- msg
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			log.Println("Consumer stopped gracefully")
			return nil
		}
	}
}

// MessageHandler adalah callback untuk handle message
type MessageHandler func(context.Context, amqp.Delivery) error

// StartForUser - consumer untuk specific user
func StartForUser(ctx context.Context, ch *amqp.Channel, userID string, handler MessageHandler) error {
	queueName := config.GetUserQueueName(userID)

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Consumer channel closed for user %s", userID)
					return
				}

				// Handle message
				if err := handler(ctx, msg); err != nil {
					log.Printf("Error handling message for user %s: %v", userID, err)
					msg.Nack(false, true) // requeue
				} else {
					msg.Ack(false)
				}

			case <-ctx.Done():
				log.Printf("Consumer stopped for user %s", userID)
				return
			}
		}
	}()

	log.Printf("✅ Consumer started for user %s", userID)
	return nil
}
