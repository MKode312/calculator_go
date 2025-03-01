package calc

import (
	"fmt"

	errors "github.com/MKode312/calculator_go/pkg/errorsForCalc"
)

func Calc(expression string) (float64, error) {
	return 0, fmt.Errorf("not implemented")
}

func Calculate(operation string, a, b float64) (float64, error) {
	switch operation {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", operation)
	}
}
