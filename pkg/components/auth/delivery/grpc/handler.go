package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/migotom/cell-centre-services/pkg/components/auth"
	"github.com/migotom/cell-centre-services/pkg/components/employee"
	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

var (
	// ContextKeyClaims defines claims context key provided by token.
	ContextKeyClaims = ContextKey("auth-claims")
	headerAuthorize  = "authorization"
)

// AuthenticateDelivery is gRPC handler delivery of authentication/authorization.
type AuthenticateDelivery struct {
	log        *zap.Logger
	repository employee.Repository
}

// NewAuthenticateDelivery returns new Auth gRPC delivery.
func NewAuthenticateDelivery(log *zap.Logger, repository employee.Repository) *AuthenticateDelivery {
	return &AuthenticateDelivery{
		log:        log,
		repository: repository,
	}
}

// DefaultInterceptor is default ahtentication/authorization interceptor, validates only token correctness without performing any role specific authorization.
func (delivery *AuthenticateDelivery) DefaultInterceptor(ctx context.Context) (context.Context, error) {
	claims, err := ObtainClaimsFromMetadata(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated with error: %v", err)
	}

	return context.WithValue(ctx, ContextKeyClaims, claims), nil
}

// Authenticate gRPC handler that checks credentials and generates AuthResponse with token and claims.
func (delivery *AuthenticateDelivery) Authenticate(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	if request == nil {
		return &pb.AuthResponse{}, grpc.Errorf(codes.Unauthenticated, "Can't authenticate: %v", auth.AuthError{Reason: auth.ErrInvalidParameters})
	}
	delivery.log.Info("Authenticate request", zap.String("request by", request.GetLogin()))

	var TokenClaimer entities.TokenClaimer

	switch e := request.GetEntity(); e {
	case pb.AuthRequest_EMPLOYEE:
		employee, err := delivery.repository.Get(context.Background(), &pb.EmployeeFilter{Email: request.GetLogin()})
		if err != nil || !auth.ValidPassword(employee.Password, request.GetPassword()) {
			return &pb.AuthResponse{}, grpc.Errorf(codes.Unauthenticated, "Can't authenticate: %v", auth.AuthError{Reason: auth.ErrInvalidCredentials})
		}
		TokenClaimer = employee

	default:
		return &pb.AuthResponse{}, grpc.Errorf(codes.Unauthenticated, "Can't authenticate: %v", auth.AuthError{Reason: auth.ErrInvalidParameters})
	}

	token, err := auth.NewToken(TokenClaimer)
	if err != nil {
		return &pb.AuthResponse{}, grpc.Errorf(codes.Unauthenticated, "Can't authenticate: %v", auth.AuthError{Reason: auth.ErrDecryptionToken, Err: err})
	}
	return &pb.AuthResponse{Token: token}, nil
}

func (delivery *AuthenticateDelivery) Validate(ctx context.Context, request *pb.ValidateRequest) (*pb.AuthResponse, error) {
	delivery.log.Debug("Validate request", zap.Any("request", request))
	return nil, nil
}

// ObtainClaimsFromMetadata obtains token claims from given context with gRPC metadata.
func ObtainClaimsFromMetadata(ctx context.Context) (claims entities.TokenClaims, err error) {
	var authenticate string
	if authenticate, err = fromMetadata(ctx); err != nil {
		return entities.TokenClaims{}, err
	}

	if claims, err = auth.ParseToken(authenticate); err != nil {
		return entities.TokenClaims{}, err
	}

	return
}

// ObtainClaimsFromContext obtains token claims from given context with value.
func ObtainClaimsFromContext(ctx context.Context) entities.TokenClaims {
	claims, ok := ctx.Value(ContextKeyClaims).(entities.TokenClaims)
	if !ok {
		return entities.TokenClaims{}
	}
	return claims
}

func fromMetadata(ctx context.Context) (authenticate string, err error) {
	if authenticate = metautils.ExtractIncoming(ctx).Get(headerAuthorize); authenticate == "" {
		return "", auth.AuthError{Reason: auth.ErrMissingToken}
	}
	return
}
