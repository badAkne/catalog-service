package rprocessor

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func reg(r *mux.Router, method, path string, handler http.Handler) {
	r.Methods(method).Path(path).Handler(handler)
}

func Loop[T any](
	ctx context.Context, wg *sync.WaitGroup,
	ch <-chan T, cb func(ctx context.Context, wg *sync.WaitGroup, obj T),
) {
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		if wg != nil {
			defer wg.Done()
		}

		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <-ch:
				if !ok {
					return
				}
				cb(ctx, wg, item)
			}
		}
	}()
}
