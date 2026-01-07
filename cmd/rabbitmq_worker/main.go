package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hrm-app/internal/pkg/rabbitmq/config"
	"hrm-app/internal/pkg/rabbitmq/connection"
	"hrm-app/internal/pkg/rabbitmq/consumer"
	"hrm-app/internal/pkg/rabbitmq/setup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go shutdown(cancel)

	conn, err := connection.New(config.RabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	if err := setup.Declare(
		ch,
		config.ExchangeName,
		config.ExchangeType,
		config.QueueName,
		config.RoutingKey,
	); err != nil {
		log.Fatal(err)
	}

	log.Println("RabbitMQ service started")
	if err := consumer.Start(ctx, ch); err != nil {
		log.Fatal(err)
	}
}

func shutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
}
