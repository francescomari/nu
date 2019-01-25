package parser

import "errors"

var (
	// ErrInvalidInput is emitted when the parser recognizes a malformed input.
	ErrInvalidInput = errors.New("invalid input")
)
