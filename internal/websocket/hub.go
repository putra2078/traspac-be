package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	// kafkautil "hrm-app/internal/pkg/kafka" // COMMENTED: Replaced with RabbitMQ
	"hrm-app/internal/middleware"
	rmqconfig "hrm-app/internal/pkg/rabbitmq/config"
	rmqconnection "hrm-app/internal/pkg/rabbitmq/connection"
	rmqmanager "hrm-app/internal/pkg/rabbitmq/manager"
	rmqpool "hrm-app/internal/pkg/rabbitmq/pool"
	rmqproducer "hrm-app/internal/pkg/rabbitmq/producer"
	rmqsetup "hrm-app/internal/pkg/rabbitmq/setup"
	"hrm-app/internal/websocket/handlerWebsocket"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	// "github.com/segmentio/kafka-go"
	// "github.com/segmentio/kafka-go" // COMMENTED: Replaced with RabbitMQ
)

// COMMENTED: Replaced with RabbitMQ
// KafkaMessage represents a message to be sent to Kafka
// type KafkaMessage struct {
// 	RoomID  uint
// 	Message []byte
// }

// RabbitMQMessage represents a message to be sent to RabbitMQ
type RabbitMQMessage struct {
	RoomID  uint
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients   map[*Client]bool
	clientsMu sync.RWMutex

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Per-user channel manager for RabbitMQ
	channelMgr  *rmqmanager.ChannelManager
	rateLimiter *middleware.RateLimiter

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Rooms for board-specific broadcasting
	rooms map[uint]map[*Client]bool

	// Chat rooms for real-time messaging
	chatRooms map[uint]map[*Client]bool

	// Redis Client
	rdb *redis.Client

	// COMMENTED: Kafka Brokers - Replaced with RabbitMQ
	// kafkaBrokers []string

	// RabbitMQ Connection and Channel Pool (for Hub-level operations)
	rabbitmqConn *amqp.Connection
	rmqPool      *rmqpool.ChannelPool
	rabbitmqURL  string

	// Instance ID for identifying the source of messages
	instanceID string

	// COMMENTED: Kafka Ingress Channel - Replaced with RabbitMQ
	// kafkaIngress chan KafkaMessage

	// RabbitMQ Ingress Channel (Worker Pool)
	rabbitmqIngress chan RabbitMQMessage

	// Global context for lifecycle management
	ctx    context.Context
	cancel context.CancelFunc

	// WaitGroup for tracking goroutines
	wg sync.WaitGroup

	// Max Clients Limit
	maxClients int

	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// NewHub creates a new Hub instance
func NewHub(rdb *redis.Client, rabbitmqURL string, channelMgr *rmqmanager.ChannelManager, rateLimiter *middleware.RateLimiter) *Hub {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize RabbitMQ connection for Hub-level operations
	conn, err := rmqconnection.New(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Create instance-specific queue name
	instanceID := uuid.NewString()
	queueName := fmt.Sprintf("chat.queue.%s", instanceID)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel for setup: %v", err)
	}

	// Declare exchange and queue
	if err := rmqsetup.Declare(ch, rmqconfig.ExchangeName, rmqconfig.ExchangeType, queueName, rmqconfig.RoutingKey); err != nil {
		log.Fatalf("Failed to declare RabbitMQ topology: %v", err)
	}
	ch.Close()

	log.Printf("RabbitMQ initialization: instance=%s, queue=%s", instanceID, queueName)

	// Initialize Channel Pool for broadcasting
	pool, err := rmqpool.NewChannelPool(conn, 10) // Pool size 10
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ channel pool: %v", err)
	}

	return &Hub{
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		channelMgr:      channelMgr,
		rateLimiter:     rateLimiter,
		shutdown:        make(chan struct{}),
		clients:         make(map[*Client]bool),
		rooms:           make(map[uint]map[*Client]bool),
		chatRooms:       make(map[uint]map[*Client]bool),
		rdb:             rdb,
		rabbitmqConn:    conn,
		rmqPool:         pool,
		rabbitmqURL:     rabbitmqURL,
		instanceID:      instanceID,
		rabbitmqIngress: make(chan RabbitMQMessage, 1000),
		ctx:             ctx,
		cancel:          cancel,
		maxClients:      10000,
	}
}

// Shutdown gracefully stops the hub and its workers
func (h *Hub) Shutdown() {
	h.shutdownOnce.Do(func() {
		log.Println("Shutting down Hub...")
		h.cancel() // Signal all goroutines to stop

		// COMMENTED: Close Kafka ingress - Replaced with RabbitMQ
		// close(h.kafkaIngress)

		// Close RabbitMQ resources
		close(h.rabbitmqIngress)
		if h.rmqPool != nil {
			h.rmqPool.Close()
		}
		if h.rabbitmqConn != nil {
			_ = h.rabbitmqConn.Close()
		}

		// Close all user channels
		log.Println("[INFO] Closing all user RabbitMQ channels...")
		// The channelManager will cleanup when connections close

		h.wg.Wait() // Wait for goroutines
		log.Println("Hub shutdown complete")
	})
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	// Start subscribing to Redis channels
	go h.subscribeToRedis()

	// COMMENTED: Start subscribing to Kafka messages - Replaced with RabbitMQ
	// go h.subscribeToKafka()

	// Start subscribing to RabbitMQ messages
	go h.subscribeToRabbitMQ()

	// COMMENTED: Start Kafka workers - Replaced with RabbitMQ
	// h.runKafkaWorkers(10) // Start 10 workers

	// Start RabbitMQ workers
	h.runRabbitMQWorkers(10) // Start 10 workers

	for {
		select {
		case <-h.ctx.Done():
			log.Println("Hub Run loop stopping due to context cancellation")
			return

		case client := <-h.register:
			h.clientsMu.Lock()
			if len(h.clients) >= h.maxClients {
				select {
				case client.send <- []byte("Error: Connection limit reached"):
				default:
				}
				close(client.send)
				h.clientsMu.Unlock()
				continue
			}
			h.clients[client] = true
			h.clientsMu.Unlock()

			log.Printf("✅ Client registered: %d (Total: %d, RabbitMQ channels: %d)",
				client.UserID, len(h.clients), h.channelMgr.GetActiveCount())

		case client := <-h.unregister:
			h.clientsMu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove client from all rooms
				for roomID, clients := range h.rooms {
					if _, ok := clients[client]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}

				// Remove client from all chat rooms
				for roomID, clients := range h.chatRooms {
					if _, ok := clients[client]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.chatRooms, roomID)
						}
					}
				}

				// Cleanup RabbitMQ resources for this user
				if client.userID != "" {
					_ = h.channelMgr.CloseUserChannel(client.userID)
					h.rateLimiter.RemoveUser(client.userID)
				} else if client.UserID != 0 {
					userIDStr := strconv.FormatUint(uint64(client.UserID), 10)
					_ = h.channelMgr.CloseUserChannel(userIDStr)
					h.rateLimiter.RemoveUser(userIDStr)
				}

				log.Printf("❌ Client unregistered: %s/%d (Total: %d, RabbitMQ channels: %d)",
					client.userID, client.UserID, len(h.clients), h.channelMgr.GetActiveCount())
			}
			h.clientsMu.Unlock()

		case message := <-h.broadcast:
			h.clientsMu.RLock()
			clientList := make([]*Client, 0, len(h.clients))
			for c := range h.clients {
				clientList = append(clientList, c)
			}
			h.clientsMu.RUnlock()

			for _, client := range clientList {
				select {
				case client.send <- message:
				default:
					// If channel blocked/closed, request unregister
					// Non-blocking send to unregister channel to prevent deadlock
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
		}
	}
}

// GetChannelManager returns the channel manager (for external access)
func (h *Hub) GetChannelManager() *rmqmanager.ChannelManager {
	return h.channelMgr
}

// GetRateLimiter returns the rate limiter (for external access)
func (h *Hub) GetRateLimiter() *middleware.RateLimiter {
	return h.rateLimiter
}

// subscribeToRedis listens for messages from Redis and forwards them to local clients
func (h *Hub) subscribeToRedis() {
	// Subscribe to all board channels
	pubsub := h.rdb.PSubscribe(h.ctx, "board:*")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-h.ctx.Done():
			log.Println("Stopping Redis subscription")
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// Extract BoardID from channel name "board:{id}"
			parts := strings.Split(msg.Channel, ":")
			if len(parts) != 2 {
				continue
			}
			boardIDStr := parts[1]
			boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
			if err != nil {
				log.Printf("Invalid board ID in redis channel: %s", msg.Channel)
				continue
			}

			h.broadcastToLocalBoard(uint(boardID), []byte(msg.Payload))
		}
	}
}

// broadcastToLocalBoard sends a message to local clients subscribed to a specific board
func (h *Hub) broadcastToLocalBoard(boardID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[boardID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				// We don't delete here to avoid locking issues or side effects in this loop
				// The client's writePump will detect error and trigger unregister
			}
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// BroadcastToBoard publishes a message to Redis channel for the board
func (h *Hub) BroadcastToBoard(boardID uint, message []byte) {
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second) // Derive from h.ctx
	defer cancel()
	channel := fmt.Sprintf("board:%d", boardID)
	if err := h.rdb.Publish(ctx, channel, message).Err(); err != nil {
		log.Printf("Error publishing to redis: %v", err)
	}
}

// BroadcastToChatRoom - Using RabbitMQ (Kafka version commented out)
func (h *Hub) BroadcastToChatRoom(roomID uint, message []byte) {
	// 1. Local broadcast dulu (immediate feedback)
	h.broadcastToLocalChatRoom(roomID, message)

	// 2. Inject InstanceID HANYA SEKALI (reuse buffer jika memungkinkan)
	enrichedMsg := h.enrichMessageWithSourceID(message)

	// COMMENTED: Kafka ingress - Replaced with RabbitMQ
	// // 3. Non-blocking queue dengan monitoring
	// select {
	// case h.kafkaIngress <- KafkaMessage{RoomID: roomID, Message: enrichedMsg}:
	// 	// Successfully queued
	// default:
	// 	// Queue full - LOG dan INCREMENT METRIC
	// 	log.Printf("[WARN] Kafka ingress queue full (size: %d), dropping message for room %d",
	// 		len(h.kafkaIngress), roomID)
	// 	// TODO: Increment dropped_messages metric
	// }

	// 3. RabbitMQ: Non-blocking queue dengan monitoring
	select {
	case h.rabbitmqIngress <- RabbitMQMessage{RoomID: roomID, Message: enrichedMsg}:
		// Successfully queued
	default:
		// Queue full - LOG dan INCREMENT METRIC
		log.Printf("[WARN] RabbitMQ ingress queue full (size: %d), dropping message for room %d",
			len(h.rabbitmqIngress), roomID)
		// TODO: Increment dropped_messages metric
	}
}

// Helper function untuk enrich message (mengurangi alokasi)
func (h *Hub) enrichMessageWithSourceID(message []byte) []byte {
	var msgData map[string]interface{}

	// Jika gagal unmarshal, kirim original message
	if err := json.Unmarshal(message, &msgData); err != nil {
		return message
	}

	// Cek apakah sudah ada _source_id (hindari duplicate)
	if _, exists := msgData["_source_id"]; exists {
		return message
	}

	// Inject source ID
	msgData["_source_id"] = h.instanceID

	enriched, err := json.Marshal(msgData)
	if err != nil {
		return message // Fallback ke original
	}

	return enriched
}

// BroadcastToChatRoomLocal just performs local broadcast
func (h *Hub) BroadcastToChatRoomLocal(roomID uint, message []byte) {
	h.broadcastToLocalChatRoom(roomID, message)
}

// RegisterClientToBoard registers a client to a specific board room
func (h *Hub) RegisterClientToBoard(clientIn handlerWebsocket.Client, boardID uint) {
	client := clientIn.(*Client)
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[boardID]; !ok {
		h.rooms[boardID] = make(map[*Client]bool)
	}
	h.rooms[boardID][client] = true
}

// RegisterClientToChatRoom registers a client to a specific chat room
func (h *Hub) RegisterClientToChatRoom(clientIn handlerWebsocket.Client, roomID uint) {
	client := clientIn.(*Client)
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.chatRooms[roomID]; !ok {
		h.chatRooms[roomID] = make(map[*Client]bool)
	}
	h.chatRooms[roomID][client] = true
}

// broadcastToLocalChatRoom - PERBAIKAN memory leak di channel
func (h *Hub) broadcastToLocalChatRoom(roomID uint, message []byte) {
	h.mu.RLock()
	clients, exists := h.chatRooms[roomID]
	h.mu.RUnlock()

	if !exists || len(clients) == 0 {
		return
	}

	// Buat client list snapshot untuk avoid holding lock
	h.mu.RLock()
	clientList := make([]*Client, 0, len(clients))
	for client := range clients {
		clientList = append(clientList, client)
	}
	h.mu.RUnlock()

	// Send ke semua clients (non-blocking)
	for _, client := range clientList {
		select {
		case client.send <- message:
			// Success
		default:
			// Channel full atau closed - schedule unregister
			go func(c *Client) {
				select {
				case h.unregister <- c:
				case <-time.After(1 * time.Second):
					// Timeout, skip
				}
			}(client)
		}
	}
}

// COMMENTED: Kafka functions - Replaced with RabbitMQ
// subscribeToKafka - PERBAIKAN dengan Health Check
// func (h *Hub) subscribeToKafka() {
// 	if len(h.kafkaBrokers) == 0 {
// 		log.Println("Kafka brokers not configured, skipping subscribeToKafka")
// 		return
// 	}
//
// 	log.Println("Starting Kafka subscription to topic: chat-messages")
//
// 	var reader *kafka.Reader
// 	consecutiveErrors := 0
// 	maxConsecutiveErrors := 10 // Batas error berturut-turut
//
// 	defer func() {
// 		if reader != nil {
// 			reader.Close()
// 			log.Println("Kafka reader closed")
// 		}
// 	}()
//
// 	for {
// 		select {
// 		case <-h.ctx.Done():
// 			log.Println("Stopping Kafka subscription loop")
// 			return
// 		default:
// 		}
//
// 		// Jika terlalu banyak error berturut-turut, pause lebih lama
// 		if consecutiveErrors >= maxConsecutiveErrors {
// 			log.Printf("[CRITICAL] Kafka consumer paused after %d consecutive errors. Waiting 30s before retry...",
// 				consecutiveErrors)
// 			time.Sleep(30 * time.Second)
// 			consecutiveErrors = 0 // Reset
// 		}
//
// 		// Create reader hanya jika belum ada
// 		if reader == nil {
// 			reader = kafkautil.GetReader(h.kafkaBrokers, "chat-messages", "hub-consumer-group")
// 			log.Println("Kafka reader created/recreated")
// 		}
//
// 		// Set deadline untuk ReadMessage using context
// 		// Ini membuat reader tidak block forever
// 		readCtx, readCancel := context.WithTimeout(h.ctx, 30*time.Second)
// 		m, err := reader.ReadMessage(readCtx)
// 		readCancel()
//
// 		if err != nil {
// 			// Check if context cancelled
// 			if h.ctx.Err() != nil {
// 				log.Println("Kafka reader context cancelled")
// 				return
// 			}
//
// 			consecutiveErrors++
//
// 			// Klasifikasi error
// 			errType := "unknown"
// 			if err == context.DeadlineExceeded {
// 				errType = "deadline_exceeded"
// 			} else if strings.Contains(err.Error(), "no data") {
// 				errType = "no_data"
// 			} else if strings.Contains(err.Error(), "connection refused") {
// 				errType = "connection_refused"
// 			}
//
// 			// Log setiap 5 error saja (kurangi spam)
// 			if consecutiveErrors%5 == 1 {
// 				log.Printf("[Kafka Consumer] Error type: %s, consecutive errors: %d, error: %v",
// 					errType, consecutiveErrors, err)
// 			}
//
// 			// Jika connection error, recreate reader
// 			if errType == "connection_refused" || consecutiveErrors >= 20 {
// 				if reader != nil {
// 					reader.Close()
// 					reader = nil
// 				}
// 				time.Sleep(5 * time.Second)
// 			} else {
// 				time.Sleep(2 * time.Second)
// 			}
//
// 			continue
// 		}
//
// 		// Successful read - reset counter
// 		consecutiveErrors = 0
//
// 		// Process message
// 		var msgData map[string]interface{}
// 		if err := json.Unmarshal(m.Value, &msgData); err == nil {
// 			// Skip jika message dari instance ini
// 			if sourceID, ok := msgData["_source_id"].(string); ok && sourceID == h.instanceID {
// 				continue
// 			}
// 		}
//
// 		// Parse RoomID
// 		roomIDStr := string(m.Key)
// 		roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
// 		if err != nil {
// 			log.Printf("Invalid room ID in Kafka message key: %s", roomIDStr)
// 			continue
// 		}
//
// 		// Broadcast ke local clients
// 		h.broadcastToLocalChatRoom(uint(roomID), m.Value)
// 	}
// }

// runKafkaWorkers starts N worker goroutines to process outgoing Kafka messages
// runKafkaWorkers - PERBAIKAN dengan Backpressure Handling
// func (h *Hub) runKafkaWorkers(numWorkers int) {
// 	for i := 0; i < numWorkers; i++ {
// 		h.wg.Add(1)
// 		go func(id int) {
// 			defer h.wg.Done()
// 			log.Printf("Kafka worker %d started", id)
//
// 			consecutiveErrors := 0
//
// 			for {
// 				select {
// 				case <-h.ctx.Done():
// 					log.Printf("Kafka worker %d stopping", id)
// 					return
//
// 				case msg, ok := <-h.kafkaIngress:
// 					if !ok {
// 						log.Printf("Kafka ingress channel closed, worker %d stopping", id)
// 						return
// 					}
//
// 					key := []byte(fmt.Sprintf("%d", msg.RoomID))
//
// 					// Gunakan context dengan timeout yang lebih masuk akal
// 					// Untuk production: 10-15 detik cukup
// 					ctx, cancel := context.WithTimeout(h.ctx, 15*time.Second)
//
// 					startTime := time.Now()
// 					err := kafkautil.ProduceMessage(ctx, key, msg.Message)
// 					duration := time.Since(startTime)
//
// 					cancel() // Always cancel
//
// 					if err != nil {
// 						consecutiveErrors++
//
// 						// Log hanya jika error pertama atau setiap 10 error
// 						if consecutiveErrors == 1 || consecutiveErrors%10 == 0 {
// 							log.Printf("[Worker %d] Kafka produce error (consecutive: %d, duration: %v): %v",
// 								id, consecutiveErrors, duration, err)
// 						}
//
// 						// Jika terlalu banyak error, pause worker sebentar
// 						if consecutiveErrors >= 5 {
// 							log.Printf("[Worker %d] Too many errors, pausing for 5s...", id)
// 							time.Sleep(5 * time.Second)
// 							consecutiveErrors = 0 // Reset setelah pause
// 						}
//
// 						// Write error to file (async, non-blocking)
// 						go h.writeKafkaError(fmt.Sprintf(
// 							"[%s] [Worker %d] Error producing to Kafka: %v (took %v)",
// 							time.Now().Format(time.RFC3339), id, err, duration))
// 					} else {
// 						// Success - reset counter
// 						consecutiveErrors = 0
//
// 						// Log slow operations
// 						if duration > 5*time.Second {
// 							log.Printf("[Worker %d] Slow Kafka produce: %v for room %d",
// 								id, duration, msg.RoomID)
// 						}
// 					}
// 				}
// 			}
// 		}(i)
// 	}
// 	log.Printf("Started %d Kafka workers", numWorkers)
// }

// Helper function to write Kafka errors (non-blocking)
func (h *Hub) writeKafkaError(errMsg string) {
	f, err := os.OpenFile("kafka_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Failed to open kafka_error.log: %v", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close kafka_error.log: %v", err)
		}
	}()

	if _, err := f.WriteString(errMsg + "\n"); err != nil {
		log.Printf("Failed to write to kafka_error.log: %v", err)
	}
}

// ========== RabbitMQ Functions (Replacing Kafka) ==========

// subscribeToRabbitMQ listens for messages from RabbitMQ and broadcasts to local clients
func (h *Hub) subscribeToRabbitMQ() {
	log.Println("Starting RabbitMQ subscription")

	// Get instance-specific queue name
	queueName := fmt.Sprintf("chat.queue.%s", h.instanceID)

	// Start consuming from the queue
	ch, err := h.rmqPool.Get()
	if err != nil {
		log.Fatalf("Failed to get channel from pool: %v", err)
	}
	defer h.rmqPool.Put(ch)

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

	log.Printf("RabbitMQ consumer started for queue: %s", queueName)

	for {
		select {
		case <-h.ctx.Done():
			log.Println("Stopping RabbitMQ subscription loop")
			return

		case msg, ok := <-msgs:
			if !ok {
				log.Println("RabbitMQ consumer channel closed")
				return
			}

			// Process message
			var msgData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &msgData); err == nil {
				// Skip if message from this instance
				if sourceID, ok := msgData["_source_id"].(string); ok && sourceID == h.instanceID {
					if err := msg.Ack(false); err != nil {
						log.Printf("Failed to ack message: %v", err)
					}
					continue
				}
			}

			// Parse RoomID from routing key or message
			// For now, we'll extract it from the message itself
			roomIDFloat, ok := msgData["room_id"].(float64)
			if !ok {
				// Try to get from headers
				if msg.Headers != nil {
					if roomIDHeader, exists := msg.Headers["room_id"]; exists {
						if roomIDInt, ok := roomIDHeader.(int64); ok {
							roomIDFloat = float64(roomIDInt)
						}
					}
				}
			}

			if roomIDFloat > 0 {
				// Broadcast to local clients
				h.broadcastToLocalChatRoom(uint(roomIDFloat), msg.Body)
			}

			// Acknowledge message
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}
}

// runRabbitMQWorkers starts N worker goroutines to process outgoing RabbitMQ messages
func (h *Hub) runRabbitMQWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		h.wg.Add(1)
		go func(id int) {
			defer h.wg.Done()
			log.Printf("RabbitMQ worker %d started", id)

			consecutiveErrors := 0

			for {
				select {
				case <-h.ctx.Done():
					log.Printf("RabbitMQ worker %d stopping", id)
					return

				case msg, ok := <-h.rabbitmqIngress:
					if !ok {
						log.Printf("RabbitMQ ingress channel closed, worker %d stopping", id)
						return
					}

					startTime := time.Now()

					// Publish to RabbitMQ using pooled channel
					ch, err := h.rmqPool.Get()
					if err != nil {
						log.Printf("[Worker %d] Failed to get channel from pool: %v", id, err)
						continue
					}

					err = rmqproducer.Publish(
						ch,
						rmqconfig.ExchangeName,
						rmqconfig.RoutingKey,
						msg.Message,
					)
					h.rmqPool.Put(ch)

					duration := time.Since(startTime)

					if err != nil {
						consecutiveErrors++

						// Log only on first error or every 10 errors
						if consecutiveErrors == 1 || consecutiveErrors%10 == 0 {
							log.Printf("[Worker %d] RabbitMQ publish error (consecutive: %d, duration: %v): %v",
								id, consecutiveErrors, duration, err)
						}

						// If too many errors, pause worker
						if consecutiveErrors >= 5 {
							log.Printf("[Worker %d] Too many errors, pausing for 5s...", id)
							time.Sleep(5 * time.Second)
							consecutiveErrors = 0 // Reset after pause
						}
					} else {
						// Success - reset counter
						consecutiveErrors = 0

						// Log slow operations
						if duration > 5*time.Second {
							log.Printf("[Worker %d] Slow RabbitMQ publish: %v for room %d",
								id, duration, msg.RoomID)
						}
					}
				}
			}
		}(i)
	}
	log.Printf("Started %d RabbitMQ workers", numWorkers)
}

// ========== End of RabbitMQ Functions ==========

// REFACTOR test
// func (h *Hub) subscribeToKafka() {
// 	if len(h.kafkaBrokers) == 0 {
// 		log.Println("Kafka brokers not configured, skipping subscribeToKafka")
// 		return
// 	}

// 	log.Println("Starting Kafka subscription to topic: chat-messages")

// 	var reader *kafka.Reader
// 	maxRetries := 5
// 	retryCount := 0

// 	// Cleanup function
// 	defer func() {
// 		if reader != nil {
// 			reader.Close()
// 			log.Println("Kafka reader closed")
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-h.ctx.Done():
// 			log.Println("Stopping Kafka subscription loop")
// 			return
// 		default:
// 		}

// 		// Create reader only if it doesn't exist
// 		if reader == nil {
// 			reader = kafkautil.GetReader(h.kafkaBrokers, "chat-messages", "hub-consumer-group")
// 			retryCount = 0
// 			log.Println("New Kafka reader created")
// 		}

// 		// Read message with parent context (no timeout here)
// 		m, err := reader.ReadMessage(h.ctx)
// 		if err != nil {
// 			// Check if context was cancelled
// 			if h.ctx.Err() != nil {
// 				log.Println("Kafka reader context cancelled")
// 				return
// 			}

// 			retryCount++
// 			errMsg := fmt.Sprintf("[%s] Error reading from Kafka (attempt %d/%d): %v",
// 				time.Now().Format(time.RFC3339), retryCount, maxRetries, err)
// 			log.Println(errMsg)

// 			// Write error to file (non-blocking)
// 			go h.writeKafkaError(errMsg)

// 			// If max retries reached, close and recreate reader
// 			if retryCount >= maxRetries {
// 				log.Println("Max retries reached, recreating Kafka reader...")
// 				if reader != nil {
// 					reader.Close()
// 					reader = nil
// 				}
// 				retryCount = 0
// 				time.Sleep(5 * time.Second) // Longer backoff
// 			} else {
// 				time.Sleep(2 * time.Second)
// 			}
// 			continue
// 		}

// 		// Reset retry count on successful read
// 		retryCount = 0

// 		// Check SourceID to avoid double delivery
// 		var msgData map[string]interface{}
// 		if err := json.Unmarshal(m.Value, &msgData); err == nil {
// 			if sourceID, ok := msgData["_source_id"].(string); ok && sourceID == h.instanceID {
// 				// Message originated from this instance, skip
// 				continue
// 			}
// 		}

// 		// Parse RoomID from message key
// 		roomIDStr := string(m.Key)
// 		roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
// 		if err != nil {
// 			log.Printf("Invalid room ID in Kafka message key: %s", roomIDStr)
// 			continue
// 		}

// 		// Broadcast to local clients
// 		h.broadcastToLocalChatRoom(uint(roomID), m.Value)
// 	}
// }

// // Helper function to write Kafka errors (non-blocking)
// func (h *Hub) writeKafkaError(errMsg string) {
// 	f, err := os.OpenFile("kafka_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Printf("Failed to open kafka_error.log: %v", err)
// 		return
// 	}
// 	defer f.Close()

// 	if _, err := f.WriteString(errMsg + "\n"); err != nil {
// 		log.Printf("Failed to write to kafka_error.log: %v", err)
// 	}
// }
