package section

import (
	"github.com/badAkne/catalog-service/internal/app/util"
)

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address        string        `required:"true"`
		Name           string        `required:"true"`
		Username       string        `required:"true"`
		Password       string        `required:"true"`
		ReadTimeout    util.Duration `default:"30s" split_words:"true"`
		WriteTimeout   util.Duration `default:"30s" split_words:"true"`
		MigrationTable string        `default:"schema_migrations" split_words:"true"`
	}
)
