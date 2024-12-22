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

	if !strings.ContainsAny(expression, "()") && len(expression) < 3 {
		isNegativeNumber := false
		if len(expression) >= 2 && expression[0] == '-' {
			isNegativeNumber = true
			for _, c := range expression[1:] {
				if c < '0' || c > '9' {
					isNegativeNumber = false
					break
				}
			}
		}
		if !isNegativeNumber {
			return 0, NewExpressionTooShortError()
		}
	}

	if len(expression) >= 3 && expression[0] == '(' && expression[len(expression)-1] == ')' {
		inner := expression[1 : len(expression)-1]
		if inner == "" {
			return 0, NewExpressionTooShortError()
		}
		if inner[0] == '-' {
			for _, c := range inner[1:] {
				if c < '0' || c > '9' {
					break
				}
			}
			return stringToFloat64(inner), nil
		}
		isNumber := true
		for _, c := range inner {
			if c < '0' || c > '9' {
				isNumber = false
				break
			}
		}
		if isNumber {
			return stringToFloat64(inner), nil
		}
	}

	for i := 0; i < len(expression); i++ {
		if expression[i] == '(' {
			openCount := 1
			j := i + 1
			for j < len(expression) && openCount > 0 {
				if expression[j] == '(' {
					openCount++
				} else if expression[j] == ')' {
					openCount--
				}
				j++
			}
			if j <= len(expression) && openCount == 0 {
				innerResult, err := Calc(expression[i+1 : j-1])
				if err != nil {
					return 0, NewBracketsExpressionError(err)
				}
				newExpr := expression[:i] + fmt.Sprintf("%g", innerResult)
				if j < len(expression) {
					newExpr += expression[j:]
				}
				return Calc(newExpr)
			}
		}
	}

	for i := 1; i < len(expression); i++ {
		if isSign(rune(expression[i])) && isSign(rune(expression[i-1])) {
			if expression[i] == '-' && (expression[i-1] == '*' || expression[i-1] == '/') {
				continue
			}
			return 0, NewConsecutiveOperatorsError()
		}
	}

	var res float64
	var b string
	var c rune = 0
	var resflag bool = false

	for _, value := range expression + "s" {
		switch {
		case value == ' ':
			continue
		case (value >= '0' && value <= '9') || value == '-':
			b += string(value)
		case isSign(value) || value == 's':
			if resflag {
				switch c {
				case '+':
					res += stringToFloat64(b)
				case '-':
					res -= stringToFloat64(b)
				case '*':
					res *= stringToFloat64(b)
				case '/':
					if stringToFloat64(b) == 0 {
						return 0, NewDivisionByZeroError()
					}
					res /= stringToFloat64(b)
				}
			} else {
				resflag = true
				res = stringToFloat64(b)
			}
			b = ""
			c = value
		default:
			return 0, NewInvalidCharError(value)
		}
	}

	return res, nil
}
