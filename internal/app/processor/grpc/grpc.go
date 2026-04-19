package grpcprocessor

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/badAkne/catalog-service/gen/grpc/catalog/v1"
	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/badAkne/catalog-service/internal/app/handler/grpc/catalog"
	"github.com/badAkne/catalog-service/internal/app/processor"
	"github.com/badAkne/catalog-service/internal/app/util"
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/grpch"
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/mcommon"
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/mprom"
	"github.com/badAkne/catalog-service/internal/pkg/http/grpc/mzerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcProc struct {
	server *grpc.Server
	addr   string
}

func NewGrpc(handler *catalog.Handler,
	unaryInterceptor []grpc.UnaryServerInterceptor,
	streamInterceptor []grpc.StreamServerInterceptor,
	cfg section.ProcessorGrpc,
) processor.Processor {
	coreMiddleware := []grpch.Middleware{
		mcommon.NewRecovery(),
		mzerolog.NewMiddleware(),
		mprom.New(),
	}

	coreUnary := make([]grpc.UnaryServerInterceptor, 0, len(coreMiddleware))
	coreStream := make([]grpc.StreamServerInterceptor, 0, len(coreMiddleware))

	for _, m := range coreMiddleware {
		coreUnary = append(coreUnary, m.ForUnary())
		coreStream = append(coreStream, m.ForStream())
	}

	unaryInterceptor = append(coreUnary, unaryInterceptor...)
	streamInterceptor = append(coreStream, streamInterceptor...)

	mprom.EnableHistogram()

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptor...),
		grpc.ChainStreamInterceptor(streamInterceptor...),
	)

	pb.RegisterCatalogServiceServer(srv, handler)
	reflection.Register(srv)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.ListenPort)

	return &grpcProc{
		server: srv,
		addr:   addr,
	}
}

func (p *grpcProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Str("listen_addr:", p.addr).Msg("Failed to start listening TCP addr for gRPC server")
		return
	}

	log.Info().Str("listen_addr", p.addr).Msg("Listening of TCP addr for gRPC server has been started")

	log.Info().Str("listener_addr", l.Addr().String()).Msg("gRPC listener address")

	go func() {
		log.Info().Msg("Starting gRPC server serve")
		err := p.server.Serve(l)
		if err != nil {
			log.Fatal().Err(err).Msg("gRPC server serve error")
		}
	}()

	processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))

	processor.WatchForShutdown(ctx, wg, util.CloserFunc(func() error {
		p.server.GracefulStop()
		return nil
	}))
}
