package authenticator

import (
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	authDelivery "github.com/migotom/cell-centre-services/pkg/components/auth/delivery/grpc"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

// Authenticator defines authenticator gRPC service.
type Authenticator struct {
	log          *zap.Logger
	config       *Config
	authDelivery *authDelivery.AuthenticateDelivery
}

// NewAuthenticator returns new authenticator gRPC service.
func NewAuthenticator(log *zap.Logger, config *Config, authDelivery *authDelivery.AuthenticateDelivery) *Authenticator {
	return &Authenticator{
		log:          log,
		config:       config,
		authDelivery: authDelivery,
	}
}

// Listen gRPC service.
func (authenticator *Authenticator) Listen() {
	listener, err := net.Listen("tcp", authenticator.config.ListenAddress)
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	var opts []grpc.ServerOption
	if authenticator.config.GRPCTLSKeyFile != "" && authenticator.config.GRPCTLSCertificateFile != "" {
		certFile := testdata.Path(authenticator.config.GRPCTLSCertificateFile)
		keyFile := testdata.Path(authenticator.config.GRPCTLSKeyFile)
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatal("Failed to generate credentials", zap.Error(err))
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, authenticator.authDelivery)
	grpcServer.Serve(listener)
}

// Config of Authenticator service.
type Config struct {
	ListenAddress          string `toml:"listen_address"`
	DatabaseAddress        string `toml:"database_address"`
	DatabaseName           string `toml:"database_name"`
	GRPCTLSCertificateFile string `toml:"grpc_tls_certificate_file"`
	GRPCTLSKeyFile         string `toml:"grpc_tls_key_file"`
}
