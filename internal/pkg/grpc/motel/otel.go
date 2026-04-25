package motel

import (
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/grpch"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type middleware struct{}

func New() grpch.Middleware { return new(middleware) }

func (*middleware) ForUnary() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor() //nolint:staticcheck
}

func (m *middleware) ForStream() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor() //nolint:staticcheck
}
