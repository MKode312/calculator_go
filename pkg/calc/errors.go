package calc

import "errors"

var (
	ErrInvalidExpression         = errors.New("invalid expression")
	ErrInvalidOp                 = errors.New("invalid operator")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrEmptyExpression           = errors.New("empty expression")
	ErrInvalidCharInExpression   = errors.New("invalid character in expression")
	ErrExpressionCannotEndWithOp = errors.New("expression cannot end with an operator")
	ErrTwoConsecutiveOps         = errors.New("two consecutive operators")
	ErrMismatchedParentheses     = errors.New("mismatched parentheses")
)
