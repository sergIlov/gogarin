package transport

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// A Handler responds to an RPC request.
//
// If ServeRPC panics, the server (the caller of ServeRPC) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, and logs a stack trace to the server error log.
type Handler interface {
	ServeRPC(ctx context.Context, req interface{}) (res interface{})
}

// NewServer constructs new RPC server.
func NewServer(conn Connection, pullTimeout time.Duration, l log.Logger) *Server {
	return &Server{
		conn:        conn,
		pollTimeout: pullTimeout,
		l:           l,
		done:        make(chan struct{}),
		m:           make(map[string]entry),
	}
}

// Server is a RPC server.
// It matches the topic of each incoming request against a list of registered
// topics and calls the corresponding handler.
//
// Server uses log polling for getting new request from a message broker.
// pollTimeout sets the limit for waiting for new messages.
// Be careful with this setting, setting it to a high value would block the Shutdown.
// The rule of thumb is to keep pollTimeout small enough for a faster Shutdown
// and large enough to not flood your message broker with a large number of requests.
type Server struct {
	conn        Connection
	pollTimeout time.Duration
	l           log.Logger

	doneMu sync.Mutex
	done   chan struct{}

	mu sync.RWMutex
	m  map[string]entry
}

type entry struct {
	h     Handler
	topic string
	done  chan struct{}
}

// ErrServerClosed is returned by the Server's Serve method after a call to Shutdown.
var ErrServerClosed = errors.New("server: Server closed")

// Handle registers the handler for the given topic.
// If a handler already exists for topic, Handle panics.
func (s *Server) Handle(topic string, handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if topic == "" {
		panic("server: invalid topic " + topic)
	}
	if handler == nil {
		panic("server: nil handler")
	}
	_, ok := s.m[topic]
	if ok {
		panic("server: multiple registrations for " + topic)
	}

	s.m[topic] = entry{h: handler, topic: topic, done: make(chan struct{})}
}

// Serve responds to incoming requests, creating a new service goroutine for each topic.
// Each service goroutine calls transport.Handler to respond to an incoming request.
func (s *Server) Serve() error {
	select {
	case <-s.done:
		return ErrServerClosed
	default:
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.serve(ctx)

	<-s.done
	cancel()
	return nil
}

func (s *Server) serve(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entry := range s.m {
		level.Info(s.l).Log("serve", entry.topic)
		go s.handle(ctx, entry.topic, entry.h, entry.done)
	}
}

func (s *Server) handle(ctx context.Context, topic string, h Handler, done chan<- struct{}) {
	var requests sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			requests.Wait()
			level.Info(s.l).Log("done", topic)
			done <- struct{}{}
			return
		default:
		}

		replyTo, data, err := s.conn.Receive(topic, s.pollTimeout)
		if err == ErrTimeout {
			continue
		}
		if err != nil {
			level.Error(s.l).Log("err", err)
			continue
		}

		requests.Add(1)
		go func() {
			defer requests.Done()
			defer func() {
				if err := recover(); err != nil {
					level.Error(s.l).Log("err", err, "serving", topic)
				}
			}()

			// TODO: Add timeout for ServeRPC
			res := h.ServeRPC(ctx, data)
			err = s.conn.Respond(replyTo, res)
			if err != nil {
				level.Error(s.l).Log("err", err)
			}
		}()
	}
}

// Shutdown gracefully shuts down the server without interrupting any active
// service goroutines. If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error.
//
// When Shutdown is called, Serve immediately returns ErrServerClosed. Make sure the
// program doesn't exit and waits instead for Shutdown to return.
func (s *Server) Shutdown(ctx context.Context) error {
	s.doneMu.Lock()
	defer s.doneMu.Unlock()

	select {
	case <-s.done:
		// Already closed. Don't close again.
		return nil
	default:
		close(s.done)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, entry := range s.m {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-entry.done:
		}
	}

	return nil
}
