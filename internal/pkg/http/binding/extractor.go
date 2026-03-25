package binding

import (
	"errors"

	"github.com/badAkne/catalog-service/internal/pkg/http/respondent"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entr "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

var (
	defaultEnTranslator ut.Translator
)
var (
	ErrMalformedSource  = errors.New("malformed request source")
	ErrValidationDailed = (*validationFailedError)(nil)
)

type validationFailedError struct {
	originalErr validator.ValidationErrors
}

func (e *validationFailedError) Error() string {
	return "Validation failed"
}

func (e *validationFailedError) Is(other error) bool {
	var err *validationFailedError

	return errors.As(other, &err)
}

func init() {
	v, _ := Validator.Engine().(*validator.Validate)

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	var found bool
	defaultEnTranslator, found = uni.GetTranslator("en")
	if !found {
		log.Panic().Msg("EN translator not found")
	}

	if err := entr.RegisterDefaultTranslations(v, defaultEnTranslator); err != nil {
		log.Panic().Msg("Failed to register EN translations: " + err.Error())
	}
}

func NewRespondentManifestExtractor(status, errorCode int, message string) respondent.ManifestExtractor {
	return func(err error) *respondent.Manifest {
		manifest := respondent.Manifest{
			Status:    status,
			ErrorCode: errorCode,
			Error:     message,
		}

		var errList validator.ValidationErrors
		if list, ok := err.(validator.ValidationErrors); ok {
			errList = list
		} else if typedErr, ok := err.(*validationFailedError); ok {
			errList = typedErr.originalErr
		} else {
			return nil
		}

		manifest.ErrorDetails = make([]string, len(errList))

		for i := 0; i < len(errList); i++ {
			if defaultEnTranslator != nil {
				manifest.ErrorDetails[i] = errList[i].Translate(defaultEnTranslator)
			} else {
				manifest.ErrorDetails[i] = errList[i].Error()
			}

		}

		return &manifest
	}
}
