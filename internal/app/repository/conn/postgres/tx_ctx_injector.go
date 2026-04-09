package rcpostgres

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type implBunIdbTxInjector struct {
	fallback bun.IDB
}

func newBunIdbTxInjector(orig bun.IDB) bun.IDB {
	return &implBunIdbTxInjector{orig}
}

func (x *implBunIdbTxInjector) getIDB(ctx context.Context) bun.IDB {
	if tx := getTxFromContext(ctx); tx.Tx != nil {
		return tx
	}

	return x.fallback
}

func (x *implBunIdbTxInjector) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	idb := x.getIDB(ctx)
	return idb.QueryContext(ctx, query, args)
}

func (x *implBunIdbTxInjector) ExecContext(
	ctx context.Context,
	query string, args ...interface{},
) (sql.Result, error) {
	idb := x.getIDB(ctx)
	return idb.ExecContext(ctx, query, args)
}

func (x *implBunIdbTxInjector) QueryRowContext(
	ctx context.Context,
	query string, args ...interface{},
) *sql.Row {
	idb := x.getIDB(ctx)
	return idb.QueryRowContext(ctx, query, args)
}

func (x *implBunIdbTxInjector) NewSelect() *bun.SelectQuery {
	query := x.fallback.NewSelect()
	return query.Conn(x)
}

func (x *implBunIdbTxInjector) NewInsert() *bun.InsertQuery {
	query := x.fallback.NewInsert()
	return query.Conn(x)
}

func (x *implBunIdbTxInjector) NewUpdate() *bun.UpdateQuery {
	query := x.fallback.NewUpdate()

	return query.Conn(x)
}

func (x *implBunIdbTxInjector) NewDelete() *bun.DeleteQuery {
	query := x.fallback.NewDelete()

	return query.Conn(x)
}

func (x *implBunIdbTxInjector) Dialect() schema.Dialect {
	return x.fallback.Dialect()
}

func (x *implBunIdbTxInjector) NewRaw(query string, args ...interface{}) *bun.RawQuery {
	return x.fallback.NewRaw(query, args).Conn(x)
}

func (x *implBunIdbTxInjector) NewValues(model interface{}) *bun.ValuesQuery {
	return x.fallback.NewValues(model).Conn(x)
}

func (x *implBunIdbTxInjector) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	return x.getIDB(ctx).BeginTx(ctx, opts)
}

func (x *implBunIdbTxInjector) RunInTx(
	ctx context.Context,
	opts *sql.TxOptions,
	f func(ctx context.Context, tx bun.Tx) error,
) error {
	return x.getIDB(ctx).RunInTx(ctx, opts, f)
}

func (x *implBunIdbTxInjector) NewCreateTable() *bun.CreateTableQuery {
	return x.fallback.NewCreateTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropTable() *bun.DropTableQuery {
	return x.fallback.NewDropTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewCreateIndex() *bun.CreateIndexQuery {
	return x.fallback.NewCreateIndex().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropIndex() *bun.DropIndexQuery {
	return x.fallback.NewDropIndex().Conn(x)
}

func (x *implBunIdbTxInjector) NewTruncateTable() *bun.TruncateTableQuery {
	return x.fallback.NewTruncateTable().Conn(x)
}

func (x *implBunIdbTxInjector) NewAddColumn() *bun.AddColumnQuery {
	return x.fallback.NewAddColumn().Conn(x)
}

func (x *implBunIdbTxInjector) NewDropColumn() *bun.DropColumnQuery {
	return x.fallback.NewDropColumn().Conn(x)
}

func (x *implBunIdbTxInjector) NewMerge() *bun.MergeQuery {
	return x.fallback.NewMerge().Conn(x)
}
