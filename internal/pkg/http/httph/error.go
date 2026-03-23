package httph

import (
	"context"
	"net/http"
)

type _ContextKeyError struct{}

type _ContextValueError struct {
	err       error
	detail    string
	isHandled bool
}

func errorPrepare(ctx context.Context) context.Context {
	var errCtx = new(_ContextValueError)
	return context.WithValue(ctx, _ContextKeyError{}, errCtx)
}

func errorGet(ctx context.Context) error {
	val := ctx.Value(_ContextKeyError{})

	if errV, ok := val.(*_ContextValueError); ok && errV != nil {
		return errV.err
	}

	return nil
}

func errorGetDetail(ctx context.Context) string {
	val := ctx.Value(_ContextKeyError{})

	if errV, ok := val.(*_ContextValueError); ok && errV != nil {
		return errV.detail
	}

	return ""
}

func errorApply(ctx context.Context, err error) {

	val := ctx.Value(_ContextKeyError{})

	if errV, ok := val.(*_ContextValueError); ok {
		errV.err = err
	}
}

func errorApplyDetail(ctx context.Context, detail string) {
	val := ctx.Value(_ContextKeyError{})

	if errV, ok := val.(*_ContextValueError); ok {
		errV.detail = detail
	}
}

func errorTryAcquireHandling(ctx context.Context) bool {
	val := ctx.Value(_ContextKeyError{})

	if errV, ok := val.(*_ContextValueError); ok && (errV == nil || errV.isHandled) {
		return false
	}

	return true
}

func ErrorPrepare(r *http.Request) *http.Request {
	return r.WithContext(errorPrepare(r.Context()))
}

func ErrorGet(r *http.Request) error {
	ctx := r.Context()

	err := errorGet(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ErrorGetDetail(r *http.Request) string {
	ctx := r.Context()

	detail := errorGetDetail(ctx)
	if detail != "" {
		return detail
	}

	return ""
}

func ErrorApply(r *http.Request, err error) {
	ctx := r.Context()

	errorApply(ctx, err)
}

func ErrorApplyDetail(r *http.Request, detail string) {
	ctx := r.Context()

	errorApplyDetail(ctx, detail)
}

func ErrorTryAcquireHandling(r *http.Request) bool {
	var isHandled bool
	ctx := r.Context()

	isHandled = errorTryAcquireHandling(ctx)

	return isHandled
}

type Middleware = func(http.Handler) http.Handler

func NewErrorMiddleware() Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, ErrorPrepare(r))
		})
	}
}
