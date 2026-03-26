package respondent

type Manifest struct {
	Status       int
	Error        string
	ErrorID      string
	ErrorCode    int
	ErrorDetail  string
	ErrorDetails []string
}

type (
	Replacer interface {
		Replace(err error) error
	}

	Expander interface {
		Expand(err error) *Manifest
	}

	Applicator interface {
		Apply(ctx any, manifest *Manifest)
	}
)
