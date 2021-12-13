package server

import (
	"context"
	"time"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/pkg/errors"

	"github.com/Octops/agones-broadcaster-http/pkg/broker"
	"github.com/Octops/agones-event-broadcaster/pkg/broadcaster"
	"k8s.io/client-go/rest"
)

type Config struct {
	Addr               string
	SyncPeriod         time.Duration
	ControllerPort     int
	MetricsBindAddress string
}

type Server struct {
	Broadcaster *broadcaster.Broadcaster
	Broker      *broker.HTTPBroker
}

func NewServer(k8sConfig *rest.Config, srvConfig *Config, broker *broker.HTTPBroker) (*Server, error) {
	gsBroadcaster := broadcaster.New(k8sConfig, broker, srvConfig.SyncPeriod, srvConfig.ControllerPort, srvConfig.MetricsBindAddress).WithWatcherFor(&v1.GameServer{})
	if gsBroadcaster == nil {
		panic(errors.New("failed to create broadcaster"))
	}

	app := &Server{
		Broadcaster: gsBroadcaster,
		Broker:      broker,
	}

	if err := gsBroadcaster.Build(); err != nil {
		panic(errors.Wrap(err, "error building broadcaster"))
	}

	return app, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.Broker.Start(ctx)

	if err := s.Broadcaster.Start(ctx); err != nil {
		return err
	}

	return nil
}
