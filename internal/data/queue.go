package data

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrQueueFull      = errors.New("data is full")
	ErrQueueNotExists = errors.New("data does not exist")
	ErrQueueLimit     = errors.New("max queues limit reached")
	ErrTimeout        = errors.New("timeout")
)

type Message struct {
	Data    string
	Created time.Time
}

type Service struct {
	queues    map[string]chan Message
	mu        sync.RWMutex
	maxQueues int
	queueCap  int
}

func New(maxQueues, queueCap int) *Service {
	return &Service{
		queues:    make(map[string]chan Message),
		maxQueues: maxQueues,
		queueCap:  queueCap,
	}
}

func (s *Service) Put(queueName string, msg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	q, exists := s.queues[queueName]
	if !exists {
		if len(s.queues) >= s.maxQueues {
			return ErrQueueLimit
		}
		q = make(chan Message, s.queueCap)
		s.queues[queueName] = q
	}

	select {
	case q <- Message{Data: msg, Created: time.Now()}:
		return nil
	default:
		return ErrQueueFull
	}
}

func (s *Service) Get(ctx context.Context, queueName string) (Message, error) {
	s.mu.RLock()
	q, exists := s.queues[queueName]
	s.mu.RUnlock()

	if !exists {
		return Message{}, ErrQueueNotExists
	}

	select {
	case msg := <-q:
		return msg, nil
	case <-ctx.Done():
		return Message{}, ErrTimeout
	}
}
