package rprocessor

import (
	"database/sql"
	"net/http"

	"github.com/badAkne/catalog-service/internal/app/entity"
	"github.com/badAkne/catalog-service/internal/pkg/http/binding"
	"github.com/badAkne/catalog-service/internal/pkg/http/httph"
	"github.com/badAkne/catalog-service/internal/pkg/http/respondent"
)

func makeErrorMiddleware() httph.Middleware {
	const (
		_40001 = "Bad request"
		_40002 = "Категория с таким названием уже существует"
		_40003 = "Товар с таким названием уже существует"
		_40401 = "Not found"
		_50001 = "Internal server error"
	)

	const (
		_40401D = "Entity was deleted or never exist"
		_50001D = "Try again later"
	)

	makeFallbackExtractor := func(status, errorCode int, message, detail string) respondent.ManifestExtractor {
		genericManifest := respondent.Manifest{
			Status:      status,
			ErrorCode:   errorCode,
			Error:       message,
			ErrorDetail: detail,
		}

		return func(_ error) *respondent.Manifest { return &genericManifest }
	}

	replacer := respondent.NewSimpleReplacer().ReplaceBy(sql.ErrNoRows, entity.ErrNotFound)

	expander := respondent.NewSimpleExpander().ExtractFor(
		binding.ErrValidationFailed,
		binding.NewRespondentManifestExtractor(http.StatusBadRequest, 40001, _40001)).
		WithoutDetail(entity.ErrCategoryAlreadyExists, http.StatusBadRequest, 40002, _40002).
		WithoutDetail(entity.ErrProductAlreadyExists, http.StatusBadRequest, 40003, _40003).
		WithDetail(entity.ErrNotFound, http.StatusNotFound, 40401, _40401, _40401D).
		FallbackExtractor(makeFallbackExtractor(http.StatusInternalServerError, 50001, _50001, _50001D))

	applicator := respondent.NewSimpleApplicator()

	return respondent.NewMiddleware(expander, replacer, applicator)
}
