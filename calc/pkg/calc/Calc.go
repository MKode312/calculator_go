package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func Apply(op rune, a, b float64) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}
	return 0, errors.New("invalid operator")
}

func Calc(expression string) (float64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	expression = strings.ReplaceAll(expression, " ", "")

	if len(expression) == 0 {
		return 0, ErrEmptyExpression
	}
	if strings.ContainsAny(string(expression[len(expression)-1]), "+-*/") {
		return 0, ErrExpressionCannotEndWithOp
	}

	operators := []rune{}
	values := []float64{}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if strings.ContainsRune("+-*/", char) {
			if i > 0 && strings.ContainsRune("+-*/", rune(expression[i-1])) {
				return 0, errors.New("two consecutive operators")
			}
		}

		if char >= '0' && char <= '9' || char == '.' {
			start := i
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, err
			}
			values = append(values, num)
			i--
		} else if char == '(' {
			operators = append(operators, char)
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				v2 := values[len(values)-1]
				values = values[:len(values)-1]
				v1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				result, err := Apply(op, v1, v2)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			if len(operators) == 0 {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		} else if _, exists := precedence[char]; exists {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				v2 := values[len(values)-1]
				values = values[:len(values)-1]
				v1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				result, err := Apply(op, v1, v2)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			operators = append(operators, char)
		} else {
			return 0, ErrInvalidCharInExpression
		}
	}

	for len(operators) > 0 {
		if len(values) < 2 {
			return 0, ErrInvalidExpression
		}
		v2 := values[len(values)-1]
		values = values[:len(values)-1]
		v1 := values[len(values)-1]
		values = values[:len(values)-1]
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]
		result, err := Apply(op, v1, v2)
		if err != nil {
			return 0, err
		}
		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}
	return values[0], nil
}
