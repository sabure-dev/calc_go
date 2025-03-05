package calc

import "errors"

var (
	ErrDivByZero       = errors.New("div by zero")
	ErrInvalidBracket  = errors.New("braces error")
	ErrInvalidOperands = errors.New("parsing error")
	ErrInvalidJson     = errors.New("json validation error")
	ErrEmptyExpression = errors.New("empty expression")

	ErrorMap = map[error]int{
		ErrInvalidBracket:  422,
		ErrInvalidOperands: 422,
		ErrDivByZero:       422,
		ErrEmptyExpression: 422,
	}
)
