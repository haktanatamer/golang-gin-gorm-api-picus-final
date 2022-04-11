package helper

import "errors"

var (
	ErrCategoryNotFound = errors.New("Category not found.")
	ErrProductNotFound  = errors.New("Product not found.")
	ErrRequestBody      = errors.New("Check your request body.")
)
