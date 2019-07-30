package event

import (
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"

	"github.com/migotom/cell-centre-backend/pkg/pb"
)

func NewEventsStreaming() *EventsStreaming {
	return &EventsStreaming{
		id: nuid.Next(),
	}
}

type EventsStreaming struct {
	sync.Mutex

	id   string
	conn stan.Conn
}

func (es *EventsStreaming) Connect(clusterID string, options ...stan.Option) (err error) {
	es.Lock()
	defer es.Unlock()

	if es.conn, err = stan.Connect(clusterID, es.id, options...); err != nil {
		return nil
	}
	return
}

// NATS returns the current NATS connection.
func (es *EventsStreaming) NATS() stan.Conn {
	es.Lock()
	defer es.Unlock()

	return es.conn
}

func (es *EventsStreaming) Publish(event *pb.Event) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if es.conn == nil {
		log.Fatal("ups")
	}
	return es.conn.Publish(event.Channel, data)
}
