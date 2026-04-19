package mcommon

import (
	"context"

	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/grpch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recoveryMiddleware struct{}

func NewRecovery() grpch.Middleware { return new(recoveryMiddleware) }

func (*recoveryMiddleware) ForUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		panicked := true

		defer func() {
			r := recover()
			if r != nil || panicked {
				err = status.Errorf(codes.Internal, "internal server error: %s", r)
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false

		return resp, err
	}
}

func (*recoveryMiddleware) ForStream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		panicked := true

		defer func() {
			r := recover()
			if r != nil || panicked {
				err = status.Errorf(codes.Internal, "internal server error: %s", r)
			}
		}()

		err = handler(srv, ss)
		panicked = false

		return err
	}
}
