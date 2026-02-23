package entity

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrNotFound              = errors.New("not found")
	ErrCategoryHasRelation   = errors.New("unable to delete category: has relation")
	ErrProductAlreadyExists  = errors.New("product already exists")
	ErrIncorrectParameters   = errors.New("incorrect parameters")
)
