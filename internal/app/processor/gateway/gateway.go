package gatewayprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/badAkne/catalog-service/gen/grpc/catalog/v1"
	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/badAkne/catalog-service/internal/app/processor"
	"github.com/badAkne/catalog-service/internal/app/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gatewayProc struct {
	server       *http.Server
	addr         string
	grpcEndpoint string
}

func NewGateway(cfg section.ProcessorGateway) processor.Processor {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.ListenPort)

	return &gatewayProc{
		addr:         addr,
		grpcEndpoint: cfg.GrpcEndpoint,
	}
}

func (g *gatewayProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	gwMux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := pb.RegisterCatalogServiceHandlerFromEndpoint(ctx, gwMux, g.grpcEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register gRPC-Gateway")
		return
	}

	g.server = &http.Server{
		Addr:              g.addr,
		Handler:           gwMux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", g.addr)
	if err != nil {
		log.Fatal().Err(err).Str("listen_addr%s", g.addr).Msg("Failed to start listening to TCP addr for gRPC-gateway")
		return
	}

	log.Info().Str("listed_addr", g.addr).Msg("Listening of TCP addr for gRPC-Gateway has been started")

	go func() {
		if err = g.server.Serve(l); err != nil {
			log.Fatal().Err(err).Msgf("Unable to start grpc server: %s", err.Error())
		}
	}()

	processor.WatchForShutdown(ctx, wg, util.CloserFunc(l.Close))

	processor.WatchForShutdown(ctx, wg, util.NewCloserContextFunc(g.server.Shutdown, ctx, 5*time.Second))
}
