package mprom

import (
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/grpch"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type middleware struct{}

func New() grpch.Middleware { return new(middleware) }

func (*middleware) ForUnary() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func (*middleware) ForStream() grpc.StreamServerInterceptor {
	return grpc_prometheus.StreamServerInterceptor
}

func EnableHistogram() {
	grpc_prometheus.EnableHandlingTimeHistogram()
}
