package calc

import "errors"

var (
	ErrInvalidExpression         = errors.New("invalid expression")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrEmptyExpression           = errors.New("empty expression")
	ErrInvalidCharInExpression   = errors.New("invalid character in expression")
	ErrExpressionCannotEndWithOp = errors.New("expression cannot end with an operator")
)
