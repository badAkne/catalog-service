package respondent

import (
	"errors"
	"net/http"

	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/rs/zerolog/log"
)

type respondent struct {
	expander   Expander
	replacer   Replacer
	applicator Applicator
}

type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
}

var ErrBadExpander = errors.New("respondent: expander is required")

func (rp *respondent) Callback(ctx any, err error) {
	newErr := rp.replacer.Replace(err)
	if errors.Is(newErr, nil) {
		return
	}

	manifest := rp.expander.Expand(newErr)
	if manifest == nil {
		return
	}

	rp.applicator.Apply(ctx, manifest)
}

func (rp *respondent) CallbackForHTTP(w http.ResponseWriter, r *http.Request, err error) {
	httpCtx := HttpContext{
		W: w,
		R: r,
	}

	rp.Callback(httpCtx, err)
}

func newRespondent(expander Expander, replacer Replacer, applicator Applicator) *respondent {
	if expander == nil {
		log.Panic().Msg(ErrBadExpander.Error())
	}

	if replacer == nil {
		replacer = NewSimpleReplacer()
	}

	if applicator == nil {
		applicator = NewSimpleApplicator()
	}

	return &respondent{
		expander:   expander,
		replacer:   replacer,
		applicator: applicator,
	}
}

func NewMiddleware(expander Expander, replacer Replacer, applicator Applicator) httph.Middleware {
	resp := newRespondent(expander, replacer, applicator)

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)

			err := httph.ErrorGet(r)

			if ok := httph.ErrorTryAcquireHandling(r); ok && !errors.Is(err, nil) {
				resp.CallbackForHTTP(w, r, err)
			}
		})
	}
}
