package mzerolog

import (
	"context"
	"time"

	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/grpch"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type middleware struct {
	log zerolog.Logger
}

func NewMiddleware() grpch.Middleware {
	return &middleware{
		log: log.Logger,
	}
}

func (m *middleware) ForUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		if err != nil {
			m.accessLogEvent(info.FullMethod, start, err)
		} else if ctx.Err() != nil {
			m.accessLogEvent(info.FullMethod, start, ctx.Err())
		}

		return resp, err
	}
}

func (m *middleware) ForStream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		ctx := ss.Context()
		start := time.Now()

		err = handler(srv, ss)

		if err != nil {
			m.accessLogEvent(info.FullMethod, start, err)
		} else if ctx.Err() != nil {
			m.accessLogEvent(info.FullMethod, start, ctx.Err())
		}

		return err
	}
}

func (m *middleware) accessLogEvent(method string, start time.Time, err error) {
	level := zerolog.TraceLevel

	if err != nil {
		level = zerolog.ErrorLevel
	}

	m.log.WithLevel(level).Err(err).Int64("duration", time.Since(start).Milliseconds()).
		Str("grpc_method", method).
		Int("grpc_status", int(status.Convert(err).Code())).
		Send()
}
