package builder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/badAkne/catalog-service/internal/app/config"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	rcategory "github.com/badAkne/catalog-service/internal/app/handler/category"
	"github.com/badAkne/catalog-service/internal/app/handler/grpc/catalog"
	rhealth "github.com/badAkne/catalog-service/internal/app/handler/health"
	rproduct "github.com/badAkne/catalog-service/internal/app/handler/product"
	"github.com/badAkne/catalog-service/internal/app/processor"
	gatewayprocessor "github.com/badAkne/catalog-service/internal/app/processor/gateway"
	grpcprocessor "github.com/badAkne/catalog-service/internal/app/processor/grpc"
	rprocessor "github.com/badAkne/catalog-service/internal/app/processor/http"
	pprocessor "github.com/badAkne/catalog-service/internal/app/processor/other"
	"github.com/badAkne/catalog-service/internal/app/repository"
	pcategory "github.com/badAkne/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/badAkne/catalog-service/internal/app/repository/product"
	rservice "github.com/badAkne/catalog-service/internal/app/service"
	mcategory "github.com/badAkne/catalog-service/internal/app/service/category"
	mproduct "github.com/badAkne/catalog-service/internal/app/service/product"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	chErrors chan error

	connPostgres *rcpostgres.Client

	categoryRepo repository.Category
	productRepo  repository.Product

	categoryService rservice.Category
	productService  rservice.Product

	healthHandler   rhandler.Health
	categoryHandler rhandler.Category
	productHandler  rhandler.Product

	grpcCatalogHandler *catalog.Handler

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	b := Builder{
		cCtx:     cCtx,
		ctx:      context.Background(),
		chErrors: make(chan error, 4096),
	}

	ctxWithCancel, cancel := context.WithCancel(b.ctx)
	b.ctx = ctxWithCancel

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go b.waitForSignal(sigChan, cancel)
	go b.printErrors()
	b.healthHandler = rhealth.NewHandler()

	return &b
}

func (b *Builder) BuildConfig(injectors ...func(c *config.Config)) {
	b.exec(true, func(b *Builder) {
		b.buildConfig(config.LoadArgs{}, injectors)
	})
}

func (b *Builder) BuildConfigSimple(injectors ...func(c *config.Config)) {
	b.exec(true, func(b *Builder) {
		b.buildConfig(config.LoadArgs{SkipConfig: true}, injectors)
	})
}

func (b *Builder) Run() {
	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
	}

	log.Info().Msg("Application is initializing")
	defer log.Info().Msg("Application is completed, GoodBye!")

	for _, proc := range b.processors {
		proc.StartAsync(b.ctx, &b.wg)
	}

	b.wg.Wait()
	log.Info().Msg("All processors finished")
}

func (b *Builder) BuildRepoConnPostgres() {
	b.exec(true, func(b *Builder) {
		cfg := b.cfg.Repository
		conn, err := rcpostgres.NewConn(b.ctx, cfg.Postgres)
		if err != nil {
			b.err = err
			return
		}

		b.connPostgres = conn
	})
}

func (b *Builder) BuildRepoConnMigrator() {
	b.exec(b.connPostgres != nil, func(b *Builder) {
		proc := pprocessor.NewMigrator(b.connPostgres)
		b.processors = append(b.processors, proc)
	})
}

func (b *Builder) BuilderRepoCategory() {
	b.exec(true, func(b *Builder) {
		categoryRepo, err := pcategory.NewRepoFromPostgres(b.ctx, b.connPostgres)
		if err != nil {
			b.err = err
			return
		}

		b.categoryRepo = categoryRepo
	}, b.connPostgres)
}

func (b *Builder) BuilderRepoProduct() {
	b.exec(true, func(b *Builder) {
		productRepo, err := pproduct.NewRepoFromPostgres(b.ctx, b.connPostgres)
		if err != nil {
			b.err = err
			return
		}

		b.productRepo = productRepo
	}, b.connPostgres)
}

func (b *Builder) exec(preCond bool, cb func(b *Builder), requiredArgs ...any) {
	if !preCond || b.err != nil {
		return
	}

	for i, reqrequiredArg := range requiredArgs {
		rv := reflect.ValueOf(reqrequiredArg)
		if !rv.IsValid() {
			b.err = fmt.Errorf("BUG: required argument #%d is nil (check dependecies)", i)
			return
		}

		if rv.Type().Kind() == reflect.Struct || !rv.IsZero() {
			continue
		}

		b.err = fmt.Errorf("BUG: required %s, but empty", rv.Type().String())
		return
	}

	cb(b)
}

func (b *Builder) buildConfig(args config.LoadArgs, injectors []func(*config.Config)) {
	args.Output = b.cCtx.App.Writer
	args.EnableSimpleLog = b.cCtx.Bool("no-json")

	config.Load(args)

	for _, injector := range injectors {
		if injector != nil {
			injector(&config.Root)
		}
	}

	b.cfg = config.Root
}

func (b *Builder) BuildServiceCategory() {
	b.exec(true, func(b *Builder) {
		service := mcategory.NewService(b.categoryRepo)
		b.categoryService = service
	}, b.categoryRepo)
}

func (b *Builder) BuildServiceProduct() {
	b.exec(true, func(b *Builder) {
		service := mproduct.NewService(b.productRepo)

		b.productService = service
	}, b.productRepo)
}

func (b *Builder) BuildHandlerHttpCategory() {
	b.exec(true, func(b *Builder) {
		handler := rcategory.NewHandler(b.categoryService)

		b.categoryHandler = handler
	}, b.categoryService)
}

func (b *Builder) BuildHandlerHttpProduct() {
	b.exec(true, func(b *Builder) {
		handler := rproduct.NewHandler(b.productService)

		b.productHandler = handler
	}, b.productService)
}

func (b *Builder) BuildHandlerGrpcCatalog() {
	b.exec(true, func(b *Builder) {
		b.grpcCatalogHandler = catalog.NewHandler(b.productService)
	}, b.productService)
}

func (b *Builder) BuildProcHttp() {
	b.exec(true, func(b *Builder) {
		procHttp := rprocessor.NewHttp(b.healthHandler, b.categoryHandler, b.productHandler, nil, b.cfg.Processor.WebServer)

		b.processors = append(b.processors, procHttp)
	}, b.healthHandler, b.productHandler, b.categoryHandler)
}

func (b *Builder) BuildProcGrpc() {
	b.exec(true, func(b *Builder) {
		procGrpc := grpcprocessor.NewGrpc(b.grpcCatalogHandler, b.cfg.Processor.Grpc)

		b.processors = append(b.processors, procGrpc)
	}, b.grpcCatalogHandler)
}

func (b *Builder) BuildProcGrpcGateway() {
	b.exec(true, func(b *Builder) {
		procGateway := gatewayprocessor.NewGateway(b.cfg.Processor.Gateway)
		b.processors = append(b.processors, procGateway)
	})
}

func (b *Builder) waitForSignal(sig chan os.Signal, cancel func()) {
	defer cancel()
	signal := <-sig
	log.Info().Msgf("Catched %s signal, shutting down", signal.String())
}

func (b *Builder) printErrors() {
	for err := range b.chErrors {
		log.Error().Err(err).Msg("Catched error from errChan")
	}
}
