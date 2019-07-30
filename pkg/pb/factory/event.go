package factory

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"

	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

const employeeEventChannel = "employees"

// EventPbFactory is pb.Event factory.
type EventPbFactory struct {
}

// NewEventPbFactory creates new factory.
func NewEventPbFactory() *EventPbFactory {
	return &EventPbFactory{}
}

// NewFromEmployeeMessage creates event based on given Employee message.
func (factory *EventPbFactory) NewFromEmployeeMessage(originator entities.TokenClaims, eventType entities.EventType, employee *pb.Employee) (*pb.Event, error) {
	event, err := newPbFromBase(originator)
	if err != nil {
		return nil, err
	}

	event.AggregateId = employee.Id
	event.Type = string(eventType)
	event.Data = &pb.Event_Employee{employee}

	return event, nil
}

// NewFromUpdateEmployeeRequest creates event based on given Employee message.
func (factory *EventPbFactory) NewFromUpdateEmployeeRequest(originator entities.TokenClaims, eventType entities.EventType, request *pb.UpdateEmployeeRequest) (*pb.Event, error) {
	event, err := newPbFromBase(originator)
	if err != nil {
		return nil, err
	}

	event.AggregateId = request.Id
	event.Type = string(eventType)
	event.Data = &pb.Event_UpdateRequest{request}

	return event, nil
}

// NewFromEmployeeFilter creates event based on given Employee message.
func (factory *EventPbFactory) NewFromEmployeeFilter(originator entities.TokenClaims, eventType entities.EventType, filter *pb.EmployeeFilter) (*pb.Event, error) {
	event, err := newPbFromBase(originator)
	if err != nil {
		return nil, err
	}

	var aggregateId string
	switch {
	case filter.GetId() != "":
		aggregateId = filter.GetId()
	case filter.GetEmail() != "":
		aggregateId = filter.GetEmail()
	}

	event.AggregateId = aggregateId
	event.Type = string(eventType)
	event.Data = &pb.Event_EmployeeFilter{filter}

	return event, nil
}

func newPbFromBase(originator entities.TokenClaims) (*pb.Event, error) {
	createdAt, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}
	return &pb.Event{
		EventId:       uuid.Must(uuid.NewV4()).String(),
		Channel:       employeeEventChannel,
		AggregateType: employeeEventChannel,
		Originator: &pb.Event_Claims{
			EntityId: originator.EntityID.Hex(),
			Entity:   originator.Entity,
			Login:    originator.Login,
		},
		CreatedAt: createdAt,
	}, nil
}
