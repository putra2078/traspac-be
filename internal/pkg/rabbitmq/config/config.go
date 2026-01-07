package config

// const (
// 	RabbitURL = "amqp://appuser:strongpassword@localhost:5672/"

// 	ExchangeName = "chat.exchange"
// 	ExchangeType = "topic"

// 	QueueName  = "jobs.queue"
// 	RoutingKey = "chat.message"

// 	WorkerCount = 5
// 	Prefetch    = 10
// )

const (
	RabbitURL    = "amqp://appuser:strongpassword@localhost:5672/"
	ExchangeName = "chat.direct" // Ubah dari "chat.exchange"
	ExchangeType = "direct"      // Ubah dari "topic"
	QueueName    = "jobs.queue"
	RoutingKey   = "chat.message"
	QueueWorker  = "chat.worker.queue"
	WorkerCount  = 5
	Prefetch     = 10
)

// GetUserQueueName returns queue name for specific user
func GetUserQueueName(userID string) string {
	return "user." + userID + ".messages"
}

// GetUserRoutingKey returns routing key for specific user (sama dengan userID)
func GetUserRoutingKey(userID string) string {
	return userID
}
