package config

import (
	"io"
	"time"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	Config struct {
		Monitor    section.Monitor
		Processor  section.Processor
		Repository section.Repository
		Meta       Meta `ignore:"true"`
	}

	Meta struct {
		WorkDir    string
		DotEnvPath string
		Load       LoadArgs
	}

	LoadArgs struct {
		Output          io.Writer `json:"-"`
		EnableSimpleLog bool
		SkipConfig      bool
	}
)

var Root Config

func Load(args LoadArgs) {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339

	if args.EnableSimpleLog {
		args.Output = zerolog.ConsoleWriter{Out: args.Output}
	}

	log.Logger = createLogger(zerolog.DebugLevel, args.Output)

	log.Debug().Msg("Logger initialized with debug level")

	if args.SkipConfig {
		log.Debug().Msg("Config loading skipped")
	}

	err := godotenv.Load()
	if err != nil {
		log.Debug().Msgf("%v\n", err)
	}

	Root.Meta.Load = args

	err = envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	level, err := zerolog.ParseLevel(Root.Monitor.LogLevel)
	if err != nil {
		log.Fatal().Msg("Unable to initialize logger")
	}

	log.Logger = createLogger(level, args.Output)
	log.Info().Msgf("Logger initialized with %s level", level)
}

func createLogger(level zerolog.Level, output io.Writer) zerolog.Logger {
	return zerolog.New(output).Level(level).With().Timestamp().Logger()
}
