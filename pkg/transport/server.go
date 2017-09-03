package transport

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
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
func NewServer(conn Connection, l log.Logger) *Server {
	return &Server{conn: conn, l: l, done: make(chan struct{}), m: make(map[string]entry)}
}

// Server is a RPC server.
// It matches the topic of each incoming request against a list of registered
// topics and calls the corresponding handler.
type Server struct {
	conn Connection
	l    log.Logger

	doneMu sync.Mutex
	done   chan struct{}

	mu sync.RWMutex
	m  map[string]entry

	workers sync.WaitGroup
}

type entry struct {
	h     Handler
	topic string
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

	s.m[topic] = entry{h: handler, topic: topic}
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
		go s.handle(ctx, entry.topic, entry.h)
	}
}

func (s *Server) handle(ctx context.Context, topic string, h Handler) {
	for {
		select {
		case <-ctx.Done():
			s.l.Log("done", topic)
			return
		default:
		}

		// TODO: Hide timeout in the connection config
		replyTo, data, err := s.conn.Receive(topic, time.Duration(2)*time.Second)
		if err == ErrTimeout {
			continue
		}
		if err != nil {
			s.l.Log("err", err)
			continue
		}

		s.workers.Add(1)
		go func() {
			defer s.workers.Done()
			defer func() {
				if err := recover(); err != nil {
					s.l.Log("err", err, "serving", topic)
				}
			}()

			// TODO: Add timeout for ServeRPC
			res := h.ServeRPC(ctx, data)
			// TODO: Add timeout for Respond
			err = s.conn.Respond(replyTo, res)
			if err != nil {
				s.l.Log("err", err)
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

	handlersDone := make(chan struct{})
	go func() {
		s.workers.Wait()
		handlersDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-handlersDone:
	}

	return nil
}
