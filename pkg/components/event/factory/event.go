package factory

import (
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson/primitive"

	employeeFactory "github.com/migotom/cell-centre-backend/pkg/components/employee/factory"
	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

// EventEntityFactory is entities.Event factory.
type EventEntityFactory struct {
	employeeFactory *employeeFactory.EmployeeEntityFactory
}

// NewEventEntityFactory creates new factory.
func NewEventEntityFactory() *EventEntityFactory {
	return &EventEntityFactory{
		employeeFactory: employeeFactory.NewEmployeeEntityFactory(nil),
	}
}

func (factory *EventEntityFactory) NewFromEvent(e pb.Event) *entities.Event {
	originatorID, _ := primitive.ObjectIDFromHex(e.Originator.EntityId)

	event := entities.Event{
		EventID:       e.EventId,
		Channel:       e.Channel,
		Type:          entities.EventType(e.Type),
		AggregateID:   e.AggregateId,
		AggregateType: e.AggregateType,
		Originator: entities.EventOriginator{
			EntityID: originatorID,
			Entity:   e.Originator.Entity,
			Login:    e.Originator.Login,
		},
	}

	switch data := e.Data.(type) {
	case *pb.Event_Employee:
		var err error
		event.Data, err = factory.employeeFactory.NewFromEmployee(data.Employee)
		if err != nil {
			return nil
		}
	case *pb.Event_UpdateRequest:
		var err error
		event.Data, err = factory.employeeFactory.NewFromUpdateEmployeeRequest(data.UpdateRequest)
		if err != nil {
			return nil
		}
	case *pb.Event_EmployeeFilter:
		var err error
		event.Data, err = factory.employeeFactory.NewFromEmployeeFilter(data.EmployeeFilter)
		if err != nil {
			return nil
		}
	}

	createdAt, err := ptypes.Timestamp(e.CreatedAt)
	if err != nil {
		return nil
	}

	event.CreatedAt = createdAt

	return &event
}
