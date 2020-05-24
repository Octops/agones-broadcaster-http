package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/Octops/agones-event-broadcaster/pkg/events"
	"github.com/sirupsen/logrus"
)

type gameserver struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Address   string            `json:"addr"`
	Port      int32             `json:"port"`
	State     string            `json:"state"`
}

type HTTPBroker struct {
	mutex sync.Mutex
	addr  string
	Store map[string]*gameserver
}

func NewHTTPBroker(addr string) *HTTPBroker {
	return &HTTPBroker{
		addr:  addr,
		Store: map[string]*gameserver{},
	}
}

func (h *HTTPBroker) Start(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.Handler))

	srv := &http.Server{
		Addr:    h.addr,
		Handler: mux,
	}

	go func() {
		logrus.Infof("server listening at %s", h.addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatal(err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				defer func() {
					cancel()
				}()

				if err := srv.Shutdown(ctxShutDown); err != nil {
					logrus.Fatalf("server shutdown failed:%+s", err)
				}
			}
		}
	}()
}

func (h *HTTPBroker) BuildEnvelope(event events.Event) (*events.Envelope, error) {
	envelope := &events.Envelope{}

	envelope.AddHeader("event_type", event.EventType().String())
	envelope.Message = event.(events.Message)

	return envelope, nil
}

func (h *HTTPBroker) SendMessage(envelope *events.Envelope) error {
	message := envelope.Message.(events.Message).Content()
	eventType := envelope.Header.Headers["event_type"]

	switch eventType {
	case "gameserver.events.added":
		gsAgones := message.(*v1.GameServer)
		return h.handleAdded(gsAgones)
	case "gameserver.events.updated":
		return h.handleUpdated(message)
	case "gameserver.events.deleted":
		gsAgones := message.(*v1.GameServer)
		return h.handleDeleted(gsAgones)
	}

	return nil
}

func (h *HTTPBroker) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.ListGameServer())
}

func (h *HTTPBroker) AddGameServer(gs *gameserver) error {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	key := fmt.Sprintf("%s/%s", gs.Namespace, gs.Name)
	h.Store[key] = gs
	return nil
}

func (h *HTTPBroker) DeleteGameServer(key string) {
	defer h.mutex.Unlock()

	h.mutex.Lock()

	delete(h.Store, key)
}

func (h *HTTPBroker) ListGameServer() map[string]*gameserver {
	defer h.mutex.Unlock()

	h.mutex.Lock()

	return h.Store
}

func (h *HTTPBroker) handleAdded(gsAgones *v1.GameServer) error {
	if gsAgones.Status.State == v1.GameServerStateReady {
		return h.AddGameServer(GameServer(gsAgones))
	}

	return nil
}

func (h *HTTPBroker) handleUpdated(message interface{}) error {
	msgUpdate := reflect.ValueOf(message)
	gsAgones := msgUpdate.Field(1).Interface().(*v1.GameServer)

	if gsAgones.Status.State == v1.GameServerStateReady {
		return h.AddGameServer(GameServer(gsAgones))
	}

	return nil
}

func (h *HTTPBroker) handleDeleted(gsAgones *v1.GameServer) error {
	key := fmt.Sprintf("%s/%s", gsAgones.Namespace, gsAgones.Name)
	h.DeleteGameServer(key)
	logrus.Infof("gameserver deleted %s", key)
	return nil
}

func GameServer(gs *v1.GameServer) *gameserver {
	return &gameserver{
		Name:      gs.Name,
		Namespace: gs.Namespace,
		Labels:    gs.Labels,
		Address:   gs.Status.Address,
		Port:      gs.Status.Ports[0].Port,
		State:     string(gs.Status.State),
	}
}