package eventlogger

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/pkg/components/event"
	eventFactory "github.com/migotom/cell-centre-services/pkg/components/event/factory"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

// EventLogger defines event logging NATS subscriber service.
type EventLogger struct {
	log             *zap.Logger
	config          *Config
	eventsStreaming *event.EventsStreaming
	eventRepository event.Repository
	eventFactory    *eventFactory.EventEntityFactory
}

// NewEventLogger returns new event logging service.
func NewEventLogger(log *zap.Logger, config *Config, eventsStreaming *event.EventsStreaming, eventRepository event.Repository) *EventLogger {
	return &EventLogger{
		log:             log,
		config:          config,
		eventsStreaming: eventsStreaming,
		eventRepository: eventRepository,
		eventFactory:    eventFactory.NewEventEntityFactory(),
	}
}

// Listen for events to log.
func (eventLogger *EventLogger) Listen() {
	streaming := eventLogger.eventsStreaming.NATS()
	if streaming == nil {
		log.Fatal("Failed to connect to NATS Streaming cluster")
	}

	for _, queue := range eventLogger.config.Subscribes {
		queueGroup := queue + "-eventlogger-group"
		queueLoggerID := "eventlogger-" + eventLogger.config.NATSClientID + "-durable"

		streaming.QueueSubscribe(queue, queueGroup, func(msg *stan.Msg) {
			var event pb.Event
			if err := proto.Unmarshal(msg.Data, &event); err != nil {
				eventLogger.log.Error("Subscribed message error", zap.String("loggerID", queueLoggerID), zap.Error(err))
				return
			}

			entity := eventLogger.eventFactory.NewFromEvent(event)
			err := eventLogger.eventRepository.New(context.Background(), entity)
			if err != nil {
				eventLogger.log.Error("Error while storing in database", zap.Error(err))
			}
		}, stan.DurableName(queueLoggerID),
		)
	}
}

// Config of EventLogger service.
type Config struct {
	DatabaseAddress string   `toml:"database_address"`
	DatabaseName    string   `toml:"database_name"`
	NATSClusterID   string   `toml:"nats_cluster_id"`
	NATSURL         string   `toml:"nats_url"`
	NATSClientID    string   `toml:"nats_client_id"`
	Subscribes      []string `toml:"subscribes"`
}
