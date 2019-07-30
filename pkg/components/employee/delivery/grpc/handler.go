package grpc

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/migotom/cell-centre-backend/pkg/components/auth"
	authDelivery "github.com/migotom/cell-centre-backend/pkg/components/auth/delivery/grpc"
	"github.com/migotom/cell-centre-backend/pkg/components/employee"
	employeeFactory "github.com/migotom/cell-centre-backend/pkg/components/employee/factory"
	"github.com/migotom/cell-centre-backend/pkg/components/event"

	"github.com/migotom/cell-centre-backend/pkg/components/role"
	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/helpers"
	"github.com/migotom/cell-centre-backend/pkg/pb"
	pbFactory "github.com/migotom/cell-centre-backend/pkg/pb/factory"
)

// TODO move into configuration
var eventsChannel = "employee"

// EmployeeDelivery is gRPC handler delivery of employee.
type EmployeeDelivery struct {
	log               *zap.Logger
	eventsStreaming   *event.EventsStreaming
	employeeFactory   *employeeFactory.EmployeeEntityFactory
	employeePbFactory *pbFactory.EmployeePbFactory
	eventPbFactory    *pbFactory.EventPbFactory
	repository        employee.Repository
}

// NewEmployeeDelivery returns new Employee gRPC delivery.
func NewEmployeeDelivery(log *zap.Logger, employeeRepository employee.Repository, roleRepository role.Repository, eventsStreaming *event.EventsStreaming) *EmployeeDelivery {
	return &EmployeeDelivery{
		log:               log,
		eventsStreaming:   eventsStreaming,
		employeeFactory:   employeeFactory.NewEmployeeEntityFactory(roleRepository),
		employeePbFactory: pbFactory.NewEmployeePbFactory(),
		eventPbFactory:    pbFactory.NewEventPbFactory(),
		repository:        employeeRepository,
	}
}

// AuthFuncOverride is authorization accessor for EmployeeDelivery gRPC with verification of allowed roles.
func (delivery *EmployeeDelivery) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	claims, err := authDelivery.ObtainClaimsFromMetadata(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated with error: %v", err)
	}

	// TODO method name based auth..
	// fmt.Println("method", fullMethodName)

	if !claims.HasRole([]string{"admin", "serviceman"}) {
		return ctx, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated with error: %v", auth.AuthError{Reason: auth.ErrInsufficientRights})
	}

	return context.WithValue(ctx, authDelivery.ContextKeyClaims, claims), nil
}

// GetEmployee gRPC handler gets employee by given filer options.
func (delivery *EmployeeDelivery) GetEmployee(ctx context.Context, filter *pb.EmployeeFilter) (*pb.Employee, error) {
	employee, err := delivery.repository.Get(context.Background(), filter)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Request unauthenticated with error: %v", err)
	}

	return delivery.employeePbFactory.NewFromEmployee(employee)
}

// NewEmployee gRPC handler creates new employee based on NewEmployeeRequest message and returns Employee message.
func (delivery *EmployeeDelivery) NewEmployee(ctx context.Context, request *pb.NewEmployeeRequest) (*pb.Employee, error) {
	if request == nil {
		return &pb.Employee{}, grpc.Errorf(codes.InvalidArgument, "Invalid request: %v", EmployeeDeliveryError{Reason: ErrInvalidEmployeeData})
	}
	if len(request.Roles) == 0 {
		return &pb.Employee{}, grpc.Errorf(codes.InvalidArgument, "Invalid request: %v", EmployeeDeliveryError{Reason: ErrInvalidEmployeeRoles})
	}

	request.Password = helpers.HashPassword(request.Password)

	employeeEntity, err := delivery.employeeFactory.NewFromNewEmployeeRequest(request)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Can't create new employee: %v", EmployeeDeliveryError{Reason: ErrInternal, Err: err})
	}
	employee, err := delivery.repository.New(context.Background(), employeeEntity)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Can't create new employee: %v", EmployeeDeliveryError{Reason: ErrInvalidEmployeeData, Err: err})
	}

	employee.Password = ""

	employeePb, err := delivery.employeePbFactory.NewFromEmployee(employee)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", EmployeeDeliveryError{Reason: ErrInternal, Err: err})
	}

	go func() {
		event, _ := delivery.eventPbFactory.NewFromEmployeeMessage(
			authDelivery.ObtainClaimsFromContext(ctx),
			entities.NewEmployeeEvent,
			employeePb,
		)
		delivery.eventsStreaming.Publish(event)
	}()
	return employeePb, nil
}

// UpdateEmployee gRPC handler updates employee data based on UpdateEmployeeRequest message and returns updated Employee message.
func (delivery *EmployeeDelivery) UpdateEmployee(ctx context.Context, request *pb.UpdateEmployeeRequest) (*pb.Employee, error) {
	if request == nil {
		return &pb.Employee{}, EmployeeDeliveryError{Reason: ErrInvalidEmployeeData}
	}

	if request.GetPassword() != "" {
		request.Password = helpers.HashPassword(request.Password)
	}

	employeeEntity, err := delivery.employeeFactory.NewFromUpdateEmployeeRequest(request)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", EmployeeDeliveryError{Reason: ErrInternal, Err: err})
	}
	employee, err := delivery.repository.Update(context.Background(), employeeEntity)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Can't update employee: %v", EmployeeDeliveryError{Reason: ErrInvalidEmployeeData, Err: err})
	}

	employeePb, err := delivery.employeePbFactory.NewFromEmployee(employee)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", EmployeeDeliveryError{Reason: ErrInternal, Err: err})
	}

	go func() {
		event, _ := delivery.eventPbFactory.NewFromUpdateEmployeeRequest(
			authDelivery.ObtainClaimsFromContext(ctx),
			entities.UpdateEmployeeEvent,
			request,
		)
		delivery.eventsStreaming.Publish(event)
	}()

	return employeePb, nil
}

// DeleteEmployee gRPC handler deletes employee based on given filter.
func (delivery *EmployeeDelivery) DeleteEmployee(ctx context.Context, filter *pb.EmployeeFilter) (*empty.Empty, error) {
	err := delivery.repository.Delete(context.Background(), filter)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Can't delete employee: %v", EmployeeDeliveryError{Reason: ErrInvalidEmployeeData, Err: err})
	}

	go func() {
		event, _ := delivery.eventPbFactory.NewFromEmployeeFilter(
			authDelivery.ObtainClaimsFromContext(ctx),
			entities.DeleteEmployeeEvent,
			filter,
		)
		delivery.eventsStreaming.Publish(event)
	}()

	return &empty.Empty{}, nil
}
