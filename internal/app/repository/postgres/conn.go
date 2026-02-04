package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewConn(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var dsn url.URL

	dsn.Scheme = "Postgres"
	dsn.Host = cfg.Address
	dsn.User = url.UserPassword(cfg.Username, cfg.Password)
	dsn.Path = cfg.Name

	args := make(url.Values)

	args.Set("sslmode", "disable")

	dsn.RawQuery = args.Encode()

	log.Printf("Write Timeout:%s\nRead Timeout:%s", cfg.WriteTimeout.String(), cfg.ReadTimeout.String())

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn.String()),
		pgdriver.WithReadTimeout(cfg.ReadTimeout.Duration),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout.Duration),
	),
	)

	sqlDB.SetMaxOpenConns(10)

	rawBunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	if err := rawBunDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}

	return &Client{
		_bunDB:   rawBunDB,
		rawBunDB: rawBunDB,
		cfg:      cfg,
	}, nil
}
