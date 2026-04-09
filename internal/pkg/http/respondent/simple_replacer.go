package respondent

import "errors"

type (
	SimpleReplacer struct {
		rules []replacementRule
	}

	replacementRule struct {
		Original    error
		Replacement error
	}
)

func (sr *SimpleReplacer) Replace(err error) error {
	for _, rule := range sr.rules {
		if errors.Is(err, rule.Original) {
			return rule.Replacement
		}
	}

	return err
}

func (sr *SimpleReplacer) ReplaceBy(original, replacement error) *SimpleReplacer {
	if !errors.Is(original, replacement) &&
		!errors.Is(original, nil) &&
		!errors.Is(replacement, nil) {
		newRule := replacementRule{
			Original:    original,
			Replacement: replacement,
		}

		sr.rules = append(sr.rules, newRule)
	}

	return sr
}

func NewSimpleReplacer() *SimpleReplacer {
	return &SimpleReplacer{
		rules: make([]replacementRule, 0),
	}
}
