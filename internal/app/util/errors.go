package util

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category doesnt exist")
	ErrCategoryHasRelation   = errors.New("unable to delete category: has relation")
	ErrProductAlreadyExists  = errors.New("product already exists")
	ErrProductNotFound       = errors.New("product doesnt exist")
)
