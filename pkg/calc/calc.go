package calc

import (
	"strconv"
	"strings"
	"unicode"
)

func Right_string(s string) bool {
	stack := []rune{}
	for _, c := range s {
		switch c {
		case ')':
			if stack[len(stack)-1] != '(' || len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		case '(':
			stack = append(stack, c)
		}
	}
	return len(stack) == 0
}

func CountOp(expression []string) bool {
	op := 0
	numbers := 0
	for _, val := range expression {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			numbers++
		} else if val == "+" || val == "-" || val == "*" || val == "/" {
			op++
		}
	}
	if numbers-op == 1 {
		return true
	} else {
		return false
	}
}

func Tokenize(expression string) []string {
	var tokens []string
	token := ""
	for _, r := range expression {
		if unicode.IsDigit(r) || r == '.' {
			token += string(r)
		} else {
			if len(token) > 0 {
				tokens = append(tokens, token)
				token = ""
			}
			tokens = append(tokens, string(r))
		}
	}
	if len(token) > 0 {
		tokens = append(tokens, token)
	}
	return tokens
}

func IsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func Calc(expression string) (float64, error) {
	if !Right_string(expression) {
		return 0.0, ErrInvalidBracket
	}
	if IsLetter(expression) {
		return 0.0, ErrInvalidOperands
	}
	expression = strings.ReplaceAll(expression, " ", "")
	tokens := Tokenize(expression)
	tokens = InfixToPostfix(tokens)
	if expression == "" || expression == " " {
		return 0.0, ErrEmptyExpression
	}
	if !CountOp(tokens) {
		return 0.0, ErrInvalidOperands
	}
	var stack []float64
	for _, val := range tokens {
		switch val {
		case "+":
			v_1 := float64(stack[len(stack)-1])
			v_2 := float64(stack[len(stack)-2])
			stack = stack[:len(stack)-2]
			stack = append(stack, float64(v_1+v_2))
		case "-":
			v_1 := float64(stack[len(stack)-1])
			v_2 := float64(stack[len(stack)-2])
			stack = stack[:len(stack)-2]
			stack = append(stack, float64(v_2-v_1))
		case "*":
			v_1 := float64(stack[len(stack)-1])
			v_2 := float64(stack[len(stack)-2])
			stack = stack[:len(stack)-2]
			stack = append(stack, float64(v_1*v_2))
		case "/":
			v_1 := float64(stack[len(stack)-2])
			v_2 := float64(stack[len(stack)-1])
			r_1 := float64(v_1)
			r_2 := float64(v_2)
			stack = stack[:len(stack)-2]
			if v_2 == 0 {
				return 0.0, ErrDivByZero
			}
			stack = append(stack, float64(r_1/r_2))
		default:
			val1, _ := strconv.ParseFloat(string(val), 64)
			stack = append(stack, float64(val1))
		}
	}
	return float64(stack[len(stack)-1]), nil
}

func InfixToPostfix(expression []string) []string {
	stack := make([]string, 0)
	postfix := []string{}
	for _, r := range expression {
		switch r {
		case "(":
			stack = append(stack, r)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

		case "+", "-":
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, r)
		case "*", "/":
			stack = append(stack, r)
		default:
			postfix = append(postfix, r)
		}
	}
	for len(stack) > 0 {
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return postfix
}
