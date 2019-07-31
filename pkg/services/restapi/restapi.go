package restapi

import (
	"context"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	gw "github.com/migotom/cell-centre-services/pkg/pb"
)

// RESTAPI defines REST API gateway service.
type RESTAPI struct {
	log    *zap.Logger
	config *Config
}

// NewRESTAPI returns new REST API gateway service.
func NewRESTAPI(log *zap.Logger, config *Config) *RESTAPI {
	return &RESTAPI{
		log:    log,
		config: config,
	}
}

// Listen to REST API calls.
func (restAPI *RESTAPI) Listen() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpc_zap.ReplaceGrpcLogger(restAPI.log)

	mux := runtime.NewServeMux()
	grpcOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(grpcOpts...),
			grpc_zap.UnaryClientInterceptor(restAPI.log, zapOpts...),
		)),
	}

	var err error
	err = gw.RegisterEmployeeServiceHandlerFromEndpoint(ctx, mux, restAPI.config.EndpointEventStoreURL, opts)
	if err != nil {
		restAPI.log.Error("Unable to register employee service handler", zap.Error(err))
		return
	}

	err = gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, restAPI.config.EndpointAuthenticatorURL, opts)
	if err != nil {
		restAPI.log.Error("Unable to register authenticator service handler", zap.Error(err))
		return
	}

	err = http.ListenAndServe(restAPI.config.ListenAddress, mux)
	if err != nil {
		restAPI.log.Error("Can't listen", zap.Error(err))
		return
	}
}

// Config of REST API service.
type Config struct {
	ListenAddress            string `toml:"listen_address"`
	EndpointEventStoreURL    string `toml:"endpoint_event_store_url"`
	EndpointAuthenticatorURL string `toml:"endpoint_authenticator_url"`
}
