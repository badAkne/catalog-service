package mzerolog

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CallbackExtractorString = func(r *http.Request) string

type CallbackExtractorAny = func(r *http.Request) any

type extractorStr struct {
	key string
	ext CallbackExtractorString
}

type extractorAny struct {
	key string
	ext CallbackExtractorAny
}

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		extStrOnSuccess []extractorStr
		extAnyOnSuccess []extractorAny

		extStrOnFail []extractorStr
		extAnyOnFail []extractorAny

		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			TailSuccess = " finished without error"
			TailFail    = " finished (or aborted) with error"
		)

		start := time.Now()

		next.ServeHTTP(w, r)

		err := httph.ErrorGet(r)

		execTime := time.Since(start)

		if m.fromOptions.skipper(r) {
			return
		}

		var mb strings.Builder
		mb.Grow(48 + len(r.RequestURI))
		mb.WriteString(r.Method)
		mb.WriteByte(' ')
		mb.WriteString(r.RequestURI)

		var ev *zerolog.Event
		var extString []extractorStr
		var extAny []extractorAny

		if errors.Is(err, nil) {
			mb.WriteString(TailSuccess)
			ev = m.log.Debug()
			extString = m.fromOptions.extStrOnSuccess
			extAny = m.fromOptions.extAnyOnSuccess
		} else {
			mb.WriteString(TailFail)
			ev = m.log.Error()
			extString = m.fromOptions.extStrOnFail
			extAny = m.fromOptions.extAnyOnFail
		}

		m.applyExtractors(r, ev, extString, extAny)

		ev.Err(err)
		ev.Str("exec_time", execTime.String())
		ev.Str("client_ip", r.RemoteAddr)

		ev.Msg(mb.String())
	})
}

func NewMiddleware(opts ...Option) httph.Middleware {
	m := middleware{
		log: log.Logger,
	}

	m.fromOptions.skipper = defaulSkipper

	for _, opt := range opts {
		opt(&m)
	}

	return m.Callback
}

func defaulSkipper(_ *http.Request) bool {
	return false
}

func (m *middleware) applyExtractors(
	r *http.Request,
	ev *zerolog.Event,
	extractorsStr []extractorStr,
	extractorsAny []extractorAny,
) {
	for _, extStr := range extractorsStr {
		key := extStr.key
		valueStr := extStr.ext(r)
		if valueStr != "" {
			ev.Str(key, valueStr)
		}
	}

	for _, extAny := range extractorsAny {
		key := extAny.key
		valueAny := extAny.ext(r)
		if valueAny != nil {
			ev.Any(key, valueAny)
		}
	}
}

/*
func newStringExtractor(key string, cb CallbackExtractorString) extractorStr {
	return extractorStr{key, cb}
}

func newAnyExtractor(key string, cb CallbackExtractorAny) extractorAny {
	return extractorAny{key, cb}
}
*/
