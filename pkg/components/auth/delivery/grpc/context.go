package grpc

// ContextKey auth gRPC context key.
type ContextKey string

func (c ContextKey) String() string {
	return "auth.delivery.grpc context key " + string(c)
}
