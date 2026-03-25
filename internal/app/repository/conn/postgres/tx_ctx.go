package rcpostgres

import (
	"context"

	"github.com/uptrace/bun"
)

type _ctxKeyTx struct{}

func getTxFromContext(ctx context.Context) bun.Tx {
	val := ctx.Value(_ctxKeyTx{})

	if tx, ok := val.(bun.Tx); ok {
		return tx
	}

	return bun.Tx{}
}

func setTxToConetext(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, _ctxKeyTx{}, tx)
}
