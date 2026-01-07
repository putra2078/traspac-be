package manager

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"hrm-app/internal/pkg/rabbitmq/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UserChannelStats struct {
	UserID           string    `json:"userId"`
	CreatedAt        time.Time `json:"createdAt"`
	LastActivity     time.Time `json:"lastActivity"`
	MessagesSent     int64     `json:"messagesSent"`
	MessagesReceived int64     `json:"messagesReceived"`
}

type UserChannel struct {
	Channel *amqp.Channel
	Stats   *UserChannelStats
	mu      sync.RWMutex
}

// GetStats returns a thread-safe copy of the user channel statistics
func (uc *UserChannel) GetStats() UserChannelStats {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return *uc.Stats
}

type ChannelManager struct {
	conn         *amqp.Connection
	userChannels sync.Map // map[string]*UserChannel
}

func NewChannelManager(conn *amqp.Connection) *ChannelManager {
	manager := &ChannelManager{
		conn: conn,
	}

	// Start idle cleanup
	go manager.startIdleCleanup()

	return manager
}

func (m *ChannelManager) CreateUserChannel(ctx context.Context, userID string) (*UserChannel, error) {
	// Check if exists
	if uc, exists := m.GetUserChannel(userID); exists {
		return uc, nil
	}

	// Create new channel
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel for user %s: %w", userID, err)
	}

	// Set QoS
	if err := ch.Qos(config.Prefetch, 0, false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	userChannel := &UserChannel{
		Channel: ch,
		Stats: &UserChannelStats{
			UserID:       userID,
			CreatedAt:    time.Now(),
			LastActivity: time.Now(),
		},
	}

	m.userChannels.Store(userID, userChannel)

	// Monitor channel
	go m.monitorChannel(userID, ch)

	log.Printf("✅ Channel created for user %s (Total: %d)", userID, m.GetActiveCount())

	return userChannel, nil
}

func (m *ChannelManager) CloseUserChannel(userID string) error {
	val, ok := m.userChannels.LoadAndDelete(userID)
	if !ok {
		return nil
	}

	uc := val.(*UserChannel)
	if err := uc.Channel.Close(); err != nil {
		log.Printf("Error closing channel for user %s: %v", userID, err)
		return err
	}

	log.Printf("❌ Channel closed for user %s (Total: %d)", userID, m.GetActiveCount())
	return nil
}

func (m *ChannelManager) GetUserChannel(userID string) (*UserChannel, bool) {
	val, ok := m.userChannels.Load(userID)
	if !ok {
		return nil, false
	}
	return val.(*UserChannel), true
}

func (m *ChannelManager) UpdateStats(userID string, statsType string) {
	uc, exists := m.GetUserChannel(userID)
	if !exists {
		return
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.Stats.LastActivity = time.Now()

	switch statsType {
	case "sent":
		uc.Stats.MessagesSent++
	case "received":
		uc.Stats.MessagesReceived++
	}
}

func (m *ChannelManager) GetActiveCount() int {
	count := 0
	m.userChannels.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (m *ChannelManager) GetAllStats() []UserChannelStats {
	stats := []UserChannelStats{}

	m.userChannels.Range(func(key, value interface{}) bool {
		uc := value.(*UserChannel)
		uc.mu.RLock()
		statsCopy := *uc.Stats
		uc.mu.RUnlock()
		stats = append(stats, statsCopy)
		return true
	})

	return stats
}

func (m *ChannelManager) monitorChannel(userID string, ch *amqp.Channel) {
	closeChan := ch.NotifyClose(make(chan *amqp.Error))

	for closeErr := range closeChan {
		if closeErr != nil {
			log.Printf("⚠️ Channel closed unexpectedly for user %s: %v", userID, closeErr)
			m.userChannels.Delete(userID)
		}
	}
}

func (m *ChannelManager) startIdleCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		idleTimeout := 30 * time.Minute

		m.userChannels.Range(func(key, value interface{}) bool {
			userID := key.(string)
			uc := value.(*UserChannel)

			uc.mu.RLock()
			lastActivity := uc.Stats.LastActivity
			uc.mu.RUnlock()

			if now.Sub(lastActivity) > idleTimeout {
				log.Printf("🧹 Closing idle channel for user %s", userID)

				if err := m.CloseUserChannel(userID); err != nil {
					log.Printf(
						"❌ Failed to close idle channel for user %s: %v",
						userID,
						err,
					)
				}
			}

			return true
		})
	}
}
