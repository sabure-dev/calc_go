package calculation

import (
	"fmt"
)

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, NewExpressionTooShortError()
	}

	bracketCount := 0
	for _, char := range expression {
		if char == '(' {
			bracketCount++
		} else if char == ')' {
			bracketCount--
		}
		if bracketCount < 0 {
			return 0, NewBracketsExpressionError(fmt.Errorf("некорректное расположение скобок"))
		}
	}
	if bracketCount != 0 {
		return 0, NewBracketsExpressionError(fmt.Errorf("некорректное расположение скобок"))
	}

	return evalExpression(expression)
}

func evalExpression(expression string) (float64, error) {
	lastOpen := -1
	for i, char := range expression {
		if char == '(' {
			lastOpen = i
		} else if char == ')' && lastOpen != -1 {
			innerResult, err := evalSimpleExpression(expression[lastOpen+1 : i])
			if err != nil {
				return 0, NewBracketsExpressionError(err)
			}

			newExpr := expression[:lastOpen]
			newExpr += fmt.Sprintf("%g", innerResult)
			newExpr += expression[i+1:]

			return evalExpression(newExpr)
		}
	}

	return evalSimpleExpression(expression)
}

func evalSimpleExpression(expression string) (float64, error) {
	var numbers []float64
	var operators []rune
	currentNumber := ""
	lastWasOperator := true

	for i, char := range expression {
		if char == ' ' {
			continue
		}

		if (char >= '0' && char <= '9') || char == '.' {
			currentNumber += string(char)
			lastWasOperator = false
			continue
		}

		if len(currentNumber) > 0 {
			num := stringToFloat64(currentNumber)
			numbers = append(numbers, num)
			currentNumber = ""
		}

		if char == '-' && (lastWasOperator || i == 0) {
			currentNumber = "-"
			lastWasOperator = true
			continue
		}

		if isOperator(char) {
			if lastWasOperator {
				return 0, NewConsecutiveOperatorsError()
			}
			operators = append(operators, char)
			lastWasOperator = true
			continue
		}

		return 0, NewInvalidCharError(char)
	}

	if len(currentNumber) > 0 {
		num := stringToFloat64(currentNumber)
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return 0, NewExpressionTooShortError()
	}

	if len(numbers) != len(operators)+1 {
		return 0, NewInvalidOperatorPositionError()
	}

	for i := 0; i < len(operators); {
		if operators[i] == '*' || operators[i] == '/' {
			if operators[i] == '*' {
				numbers[i+1] = numbers[i] * numbers[i+1]
			} else {
				if numbers[i+1] == 0 {
					return 0, NewDivisionByZeroError()
				}
				numbers[i+1] = numbers[i] / numbers[i+1]
			}
			numbers = append(numbers[:i], numbers[i+1:]...)
			operators = append(operators[:i], operators[i+1:]...)
		} else {
			i++
		}
	}

	result := numbers[0]
	for i := 0; i < len(operators); i++ {
		if operators[i] == '+' {
			result += numbers[i+1]
		} else if operators[i] == '-' {
			result -= numbers[i+1]
		}
	}

	return result, nil
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func stringToFloat64(s string) float64 {
	var result float64
	isNegative := false

	if s == "" {
		return 0
	}

	if s[0] == '-' {
		isNegative = true
		s = s[1:]
	}

	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + float64(c-'0')
		}
	}

	if isNegative {
		result = -result
	}

	return result
}
