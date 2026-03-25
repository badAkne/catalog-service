package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/badAkne/catalog-service/internal/app/config/section"
	"github.com/badAkne/catalog-service/migration"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
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

	bunDB := newBunIdbTxInjector(rawBunDB)

	return &Client{
		_bunDB:   bunDB,
		rawBunDB: rawBunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()

	err = migrations.Discover(migration.Postgres)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to get migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	m := migrate.NewMigrator(c.GetRawBunDB(), migrations, opts...)

	err = m.Init(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to init migrations: %w", err)
	}

	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to get applied migrations:%w", err)
	}

	if len(applied) > 0 {
		oldVer, err = strconv.ParseInt(applied[0].Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("unable to parse int: %w", err)
		}
	}

	mgg, err := m.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to migrate: %w", err)
	}

	if !mgg.IsZero() {
		for _, mig := range mgg.Migrations {
			newVer = max(newVer, mig.ID)
		}
	} else {
		newVer = oldVer
	}

	return oldVer, newVer, nil
}

func (c *Client) InsideTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx := getTxFromContext(ctx)
	if tx.Tx != nil {
		return f(ctx)
	}

	var done bool
	var err error

	tx, err = c.rawBunDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	ctx = setTxToConetext(ctx, tx)

	err = f(ctx)
	if err != nil {
		return err
	}

	done = true
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
