package eventstore

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	authDelivery "github.com/migotom/cell-centre-services/pkg/components/auth/delivery/grpc"
	"github.com/migotom/cell-centre-services/pkg/components/employee"
	employeeDelivery "github.com/migotom/cell-centre-services/pkg/components/employee/delivery/grpc"
	"github.com/migotom/cell-centre-services/pkg/components/event"
	"github.com/migotom/cell-centre-services/pkg/components/role"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

// EventStore defines gRPC handler for REST API and NATS events propagator.
type EventStore struct {
	log              *zap.Logger
	config           *Config
	authDelivery     *authDelivery.AuthenticateDelivery
	employeeDelivery *employeeDelivery.EmployeeDelivery
}

// NewEventStore returns new event store service.
func NewEventStore(
	log *zap.Logger,
	config *Config,
	eventsStreaming *event.EventsStreaming,
	authDelivery *authDelivery.AuthenticateDelivery,
	employeeRepository employee.Repository,
	roleRepository role.Repository,
) *EventStore {
	return &EventStore{
		log:              log,
		config:           config,
		authDelivery:     authDelivery,
		employeeDelivery: employeeDelivery.NewEmployeeDelivery(log, employeeRepository, roleRepository, eventsStreaming),
	}
}

// Listen to gRPC calls.
func (eventStore *EventStore) Listen() {
	listener, err := net.Listen("tcp", eventStore.config.ListenAddress)
	if err != nil {
		eventStore.log.Fatal("Failed to listen", zap.Error(err))
	}

	var opts []grpc.ServerOption
	if eventStore.config.GRPCTLSKeyFile != "" && eventStore.config.GRPCTLSCertificateFile != "" {
		certFile := testdata.Path(eventStore.config.GRPCTLSCertificateFile)
		keyFile := testdata.Path(eventStore.config.GRPCTLSKeyFile)
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			eventStore.log.Fatal("Failed to generate credentials", zap.Error(err))
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_auth.UnaryServerInterceptor(eventStore.authDelivery.DefaultInterceptor),
	)))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(grpcServer, eventStore.employeeDelivery)
	grpcServer.Serve(listener)
}

// Config of EventStore service.
type Config struct {
	ListenAddress          string `toml:"listen_address"`
	DatabaseAddress        string `toml:"database_address"`
	DatabaseName           string `toml:"database_name"`
	NATSClusterID          string `toml:"nats_cluster_id"`
	NATSURL                string `toml:"nats_url"`
	GRPCTLSCertificateFile string `toml:"grpc_tls_certificate_file"`
	GRPCTLSKeyFile         string `toml:"grpc_tls_key_file"`
}
