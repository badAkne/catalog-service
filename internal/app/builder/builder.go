package builder

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/badAkne/catalog-service/internal/app/config"
	rhandler "github.com/badAkne/catalog-service/internal/app/handler"
	rhealth "github.com/badAkne/catalog-service/internal/app/handler/health"
	"github.com/badAkne/catalog-service/internal/app/processor"
	pprocessor "github.com/badAkne/catalog-service/internal/app/processor/other"
	"github.com/badAkne/catalog-service/internal/app/repository"
	pcategory "github.com/badAkne/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/badAkne/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/badAkne/catalog-service/internal/app/repository/product"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	connPostgres *rcpostgres.Client

	categoryRepo repository.Category
	productRepo  repository.Product

	healthHandler rhandler.Health
	//TODO: добавить обратно, линтер ругается
	//categoryHandler rhandler.Category
	//productHandler  rhandler.Product

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	var b = Builder{
		cCtx: cCtx,
		ctx:  context.Background(),
	}

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
