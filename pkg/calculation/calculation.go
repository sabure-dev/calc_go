package calculation

import (
	"fmt"
	"strings"
)

func stringToFloat64(str string) float64 {
	if str == "" {
		return 0
	}

	isNegative := str[0] == '-'
	if isNegative {
		str = str[1:]
	}

	var result float64
	for _, digit := range str {
		if digit >= '0' && digit <= '9' {
			result = result*10 + float64(digit-'0')
		}
	}

	if isNegative {
		result = -result
	}

	return result
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, NewExpressionTooShortError()
	}

	// Проверка скобок
	openBrackets := 0
	for _, char := range expression {
		if char == '(' {
			openBrackets++
		} else if char == ')' {
			openBrackets--
			if openBrackets < 0 {
				return 0, NewBracketsExpressionError(fmt.Errorf("лишняя закрывающая скобка"))
			}
		}
	}
	if openBrackets > 0 {
		return 0, NewBracketsExpressionError(fmt.Errorf("не хватает закрывающей скобки"))
	}

	// Обработка скобок
	bracketStart := -1
	maxDepth := 0
	currentDepth := 0

	// Находим самые глубокие скобки
	for i, char := range expression {
		if char == '(' {
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
				bracketStart = i
			}
		} else if char == ')' && currentDepth == maxDepth && bracketStart != -1 {
			// Вычисляем значение в скобках
			innerResult, err := Calc(expression[bracketStart+1 : i])
			if err != nil {
				return 0, err
			}

			// Заменяем скобочное выражение на результат
			innerStr := fmt.Sprintf("%g", innerResult)
			if innerResult < 0 {
				innerStr = "(" + innerStr + ")"
			}
			expression = expression[:bracketStart] + innerStr + expression[i+1:]

			// Сбрасываем поиск скобок
			bracketStart = -1
			maxDepth = 0
			currentDepth = 0
			i = -1 // При следующей итерации станет 0
			continue
		} else if char == ')' {
			currentDepth--
		}
	}

	// Проверка на простое отрицательное число
	if len(expression) >= 2 && expression[0] == '-' {
		isNumber := true
		for _, c := range expression[1:] {
			if c < '0' || c > '9' {
				isNumber = false
				break
			}
		}
		if isNumber {
			return stringToFloat64(expression), nil
		}
	}

	if !strings.ContainsAny(expression, "+-*/") {
		if len(expression) == 0 {
			return 0, NewExpressionTooShortError()
		}
		return stringToFloat64(expression), nil
	}

	var numbers []float64
	var operators []rune
	var currentNumber string
	var lastChar rune
	var lastWasOperator bool

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		switch {
		case char == ' ':
			continue
		case char >= '0' && char <= '9':
			currentNumber += string(char)
			lastWasOperator = false
		case char == '-' && (i == 0 || lastChar == '(' || isSign(lastChar)):
			if lastWasOperator && lastChar != '(' {
				return 0, NewConsecutiveOperatorsError()
			}
			currentNumber = string(char)
			lastWasOperator = false
		case isSign(char):
			if lastWasOperator {
				return 0, NewConsecutiveOperatorsError()
			}
			if currentNumber != "" {
				numbers = append(numbers, stringToFloat64(currentNumber))
				currentNumber = ""
			}
			operators = append(operators, char)
			lastWasOperator = true
		case char == '(' || char == ')':
			if currentNumber != "" {
				numbers = append(numbers, stringToFloat64(currentNumber))
				currentNumber = ""
			}
			lastWasOperator = false
		default:
			return 0, NewInvalidCharError(char)
		}
		lastChar = char
	}

	if currentNumber != "" {
		numbers = append(numbers, stringToFloat64(currentNumber))
	}

	if len(numbers) == 0 || len(numbers) != len(operators)+1 {
		return 0, NewInvalidOperatorPositionError()
	}

	// Вычисление умножения и деления
	for i := 0; i < len(operators); i++ {
		if operators[i] == '*' || operators[i] == '/' {
			if operators[i] == '*' {
				numbers[i+1] = numbers[i] * numbers[i+1]
			} else {
				if numbers[i+1] == 0 {
					return 0, NewDivisionByZeroError()
				}
				numbers[i+1] = numbers[i] / numbers[i+1]
			}
			numbers[i] = 0
			operators[i] = '+'
		}
	}

	// Вычисление сложения и вычитания
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
