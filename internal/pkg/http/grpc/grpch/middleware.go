package grpch

import "google.golang.org/grpc"

type Middleware interface {
	ForUnary() grpc.UnaryServerInterceptor
	ForStream() grpc.StreamServerInterceptor
}
