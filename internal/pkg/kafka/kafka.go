package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	Writer *kafka.Writer
)

func InitKafka(brokers []string) {
	if Writer != nil {
		return
	}

	Writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        "chat-messages",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    10,                    // Lebih kecil untuk chat
		BatchTimeout: 100 * time.Millisecond, // Lebih lama
		ReadTimeout:  30 * time.Second,      // Lebih lama
		WriteTimeout: 30 * time.Second,      // Lebih lama
		MaxAttempts:  3,
		RequiredAcks: kafka.RequireOne,      // Tambahkan ini
		Async:        false,                  // Synchronous untuk reliability
	}
	log.Println("Kafka Writer initialized")
}

func GetReader(brokers []string, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         groupID,
		Topic:           topic,
		MinBytes:        1,                    // ✅ 1 byte minimum
		MaxBytes:        1e6,                  // 1MB max
		CommitInterval:  time.Second,          // Auto-commit setiap detik
		StartOffset:     kafka.LastOffset,     // Mulai dari offset terakhir
		MaxWait:         500 * time.Millisecond, // Max wait untuk batch
		ReadBackoffMin:  100 * time.Millisecond,
		ReadBackoffMax:  1 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    10 * time.Second,
	})
}

func ProduceMessage(ctx context.Context, key, value []byte) error {
	if Writer == nil {
		return fmt.Errorf("kafka writer not initialized")
	}
	
	return Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
}

func CloseKafka() {
	if Writer != nil {
		if err := Writer.Close(); err != nil {
			log.Printf("failed to close writer: %v", err)
		} else {
			log.Println("Kafka Writer closed successfully")
		}
	}
}