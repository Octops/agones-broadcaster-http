package server

import (
	"context"

	"github.com/Octops/agones-broadcaster-http/pkg/broker"
	"github.com/Octops/agones-event-broadcaster/pkg/broadcaster"
	"k8s.io/client-go/rest"
)

type Server struct {
	Broadcaster *broadcaster.Broadcaster
	Broker      *broker.HTTPBroker
}

func NewServer(config *rest.Config, addr string) (*Server, error) {
	httpBroker := broker.NewHTTPBroker(addr)
	gsBroadcaster, err := broadcaster.New(config, httpBroker)
	if err != nil {
		return nil, err
	}

	app := &Server{
		Broadcaster: gsBroadcaster,
		Broker:      httpBroker,
	}

	return app, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.Broker.Start(ctx)

	if err := s.Broadcaster.Start(); err != nil {
		return err
	}

	return nil
}
