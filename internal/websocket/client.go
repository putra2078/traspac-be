/*
package websocket

import (
	"context"
	"encoding/json" // Added for marshaling error messages
	"log"
	"sync" // Added for RateLimiter
	"time"



	"github.com/gorilla/websocket"
)

// RateLimiter implements a simple Token Bucket rate limiter
type RateLimiter struct {
	rate       float64   // tokens per second
	capacity   float64   // max tokens
	tokens     float64   // current tokens
	lastRefill time.Time // last time tokens were added
	mu         sync.Mutex
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // start full
		lastRefill: time.Now(),
	}
}

// IsAllowed checks if the action is allowed, consuming 1 token if so
func (rl *RateLimiter) IsAllowed() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	// Refill tokens
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastRefill = now

	// Consume token
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	return false
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// UserID of the connected client
	UserID uint

	// UserName of the connected client
	UserName string

	// UserUsername of the connected client
	UserUsername string

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Context for canceling the connection
	Ctx    context.Context
	Cancel context.CancelFunc

	// Rate Limiter
	limiter *RateLimiter
}

func (c *Client) GetUserID() uint {
	return c.UserID
}

func (c *Client) GetUserName() string {
	return c.UserName
}

func (c *Client) GetUserUsername() string {
	return c.UserUsername
}

func (c *Client) GetContext() context.Context {
	return c.Ctx
}

func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		// Channel is full or closed, handlePump will eventually clean up
	}
}

func (c *Client) Close() {
	if c.Cancel != nil {
		c.Cancel()
	}
	c.conn.Close()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		if c.Cancel != nil {
			c.Cancel()
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Rate Limit Check
		if c.limiter != nil && !c.limiter.IsAllowed() {
			log.Printf("Rate limit exceeded for user %d", c.UserID)
			// Send error message
			errorMsg := map[string]string{"action": "error", "message": "Rate limit exceeded. Connection closed."}
			if jsonMsg, err := json.Marshal(errorMsg); err == nil {
				c.conn.WriteMessage(websocket.TextMessage, jsonMsg)
			}
			// Close connection immediately
			// The defer block will handle unregister/cleanup
			return
		}

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.Ctx.Done():
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
*/

package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"

	rmqConfig "hrm-app/internal/pkg/rabbitmq/config"
	"hrm-app/internal/pkg/rabbitmq/consumer"
	"hrm-app/internal/pkg/rabbitmq/producer"
	"hrm-app/internal/pkg/rabbitmq/setup"
)

// RateLimiter implements a simple Token Bucket rate limiter (maintained for backward compatibility with handler.go)
type RateLimiter struct {
	rate       float64   // tokens per second
	capacity   float64   // max tokens
	tokens     float64   // current tokens
	lastRefill time.Time // last time tokens were added
	mu         sync.Mutex
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // start full
		lastRefill: time.Now(),
	}
}

// IsAllowed checks if the action is allowed, consuming 1 token if so
func (rl *RateLimiter) IsAllowed() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	// Refill tokens
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastRefill = now

	// Consume token
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	return false
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512 KB
)

type Client struct {
	hub          *Hub
	conn         *websocket.Conn
	send         chan []byte
	userID       string
	UserID       uint
	UserName     string
	UserUsername string
	Ctx          context.Context
	Cancel       context.CancelFunc
	limiter      *RateLimiter
}

// Compatibility methods for handlerWebsocket.Client interface
func (c *Client) GetUserID() uint         { return c.UserID }
func (c *Client) GetUserName() string     { return c.UserName }
func (c *Client) GetUserUsername() string { return c.UserUsername }
func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
	}
}
func (c *Client) Close() {
	if c.Cancel != nil {
		c.Cancel()
	}
	_ = c.conn.Close()
}
func (c *Client) GetContext() context.Context { return c.Ctx }

type Message struct {
	Type      string          `json:"type"`
	To        string          `json:"to,omitempty"`
	From      string          `json:"from,omitempty"`
	Content   string          `json:"content,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		Ctx:    ctx,
		Cancel: cancel,
	}

	// Register client
	client.hub.register <- client

	// Setup RabbitMQ for this user
	if err := client.setupRabbitMQ(); err != nil {
		log.Printf("Failed to setup RabbitMQ for %s: %v", userID, err)
		_ = client.conn.Close()
		return
	}

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

func (c *Client) setupRabbitMQ() error {
	channelMgr := c.hub.GetChannelManager()

	// Create user channel
	uc, err := channelMgr.CreateUserChannel(c.Ctx, c.userID)
	if err != nil {
		return err
	}

	// Setup topology
	if err := setup.DeclareUserQueue(uc.Channel, rmqConfig.ExchangeName, rmqConfig.ExchangeType, c.userID); err != nil {
		return err
	}

	if err := consumer.StartForUser(c.Ctx, uc.Channel, c.userID, c.handleRabbitMQMessage); err != nil {
		return err
	}

	log.Printf("✅ RabbitMQ setup complete for user %s", c.userID)
	return nil
}

func (c *Client) handleRabbitMQMessage(ctx context.Context, msg amqp.Delivery) error {
	// Update stats
	c.hub.channelMgr.UpdateStats(c.userID, "received")

	// Forward to WebSocket
	select {
	case c.send <- msg.Body:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Cancel()
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for %s: %v", c.userID, err)
			}
			break
		}

		// Rate limiting
		rateLimiter := c.hub.GetRateLimiter()
		result := rateLimiter.CheckLimit(c.userID, 10, 60) // 10 msg per 60 sec

		if !result.Allowed {
			errorMsg := Message{
				Type:    "error",
				Content: "Rate limit exceeded. Please slow down.",
				Data:    json.RawMessage(`{"retryAfter":` + string(rune(result.RetryAfter)) + `}`),
			}
			if errBytes, err := json.Marshal(errorMsg); err == nil {
				c.send <- errBytes
			}
			log.Printf("⚠️ Rate limit exceeded for user %s", c.userID)
			continue
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Invalid message from %s: %v", c.userID, err)
			continue
		}

		msg.From = c.userID
		msg.Timestamp = time.Now().Unix()

		// Handle based on type
		if err := c.handleMessage(&msg); err != nil {
			log.Printf("Error handling message from %s: %v", c.userID, err)
		}
	}
}

func (c *Client) handleMessage(msg *Message) error {
	channelMgr := c.hub.GetChannelManager()
	uc, exists := channelMgr.GetUserChannel(c.userID)
	if !exists {
		return nil
	}

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	switch msg.Type {
	case "chat":
		// Send to specific user
		if msg.To != "" {
			if err := producer.PublishToUser(uc.Channel, rmqConfig.ExchangeName, msg.To, messageBytes, c.userID); err != nil {
				return err
			}
			channelMgr.UpdateStats(c.userID, "sent")
			log.Printf("📤 %s -> %s: %s", c.userID, msg.To, msg.Content)
		}

	case "broadcast":
		// Broadcast to all users (optional)
		c.hub.broadcast <- messageBytes

	default:
		log.Printf("Unknown message type: %s from %s", msg.Type, c.userID)
	}

	return nil
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.Ctx.Done():
			return
		}
	}
}
