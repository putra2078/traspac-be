package pool

import (
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelPool struct {
	conn       *amqp.Connection
	channels   chan *amqp.Channel
	sharedPool []*amqp.Channel
	size       int
	mu         sync.Mutex
}

func NewChannelPool(conn *amqp.Connection, size int) (*ChannelPool, error) {
	pool := &ChannelPool{
		conn:       conn,
		channels:   make(chan *amqp.Channel, size),
		sharedPool: make([]*amqp.Channel, size),
		size:       size,
	}

	for i := 0; i < size; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to open channel for pool: %w", err)
		}
		pool.channels <- ch
		pool.sharedPool[i] = ch
	}

	log.Printf("🌊 Channel pool initialized with size %d", size)
	return pool, nil
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	select {
	case ch := <-p.channels:
		if ch.IsClosed() {
			newCh, err := p.conn.Channel()
			if err != nil {
				return nil, err
			}
			return newCh, nil
		}
		return ch, nil
	default:
		// If pool empty (shouldn't happen if sized correctly), create temporary channel
		return p.conn.Channel()
	}
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	if ch == nil || ch.IsClosed() {
		return
	}

	select {
	case p.channels <- ch:
		// Returned to pool
	default:
		// Pool full, close channel
		_ = ch.Close()
	}
}

// GetShared returns a channel from the shared pool based on a string ID (sticky)
func (p *ChannelPool) GetShared(id string) *amqp.Channel {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Simple hash for index
	var hash int
	for i := 0; i < len(id); i++ {
		hash = hash*31 + int(id[i])
	}

	index := hash % p.size
	if index < 0 {
		index = -index
	}

	ch := p.sharedPool[index]
	if ch.IsClosed() {
		// Recreate channel if closed
		newCh, err := p.conn.Channel()
		if err == nil {
			p.sharedPool[index] = newCh
			return newCh
		}
	}
	return ch
}

func (p *ChannelPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.channels)
	for ch := range p.channels {
		_ = ch.Close()
	}
	log.Println("🌊 Channel pool closed")
}
