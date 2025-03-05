package core

import (
	"github.com/sabure-dev/calc_go/internal/agent/core"
	"github.com/sabure-dev/calc_go/pkg/calc"
	"strconv"
	"strings"
)

type node struct {
	value  string
	left   *node
	right  *node
	result float64
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func buildAST(tokens []string) *node {
	var stack []*node
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, &node{value: token, left: left, right: right})
		default:
			stack = append(stack, &node{value: token})
		}
	}
	return stack[0]
}

func calcLvl(node *node, tasks chan core.Task, results chan float64, id int) (float64, error) {
	if !isOperator(node.value) {
		val, _ := strconv.ParseFloat(node.value, 64)
		return val, nil
	}
	left, err := calcLvl(node.left, tasks, results, id)
	if err != nil {
		return 0, err
	}
	right, err := calcLvl(node.right, tasks, results, id)
	if err != nil {
		return 0, err
	}
	if node.value == "/" && right == 0 {
		return 0.0, calc.ErrDivByZero
	}
	task := core.Task{Operation: node.value, Arg1: left, Arg2: right, Result: node.result, Id: id}
	tasks <- task
	result := <-results
	return result, nil
}

func Calc(expression string, tasks chan core.Task, results chan float64, errors chan error, done chan int, id int) {
	if !calc.Right_string(expression) {
		errors <- calc.ErrInvalidBracket
		done <- 1
		return
	}
	if calc.IsLetter(expression) {
		errors <- calc.ErrInvalidOperands
		done <- 1
		return
	}
	if expression == "" || expression == " " {
		errors <- calc.ErrEmptyExpression
		done <- 1
		return
	}
	expression = strings.ReplaceAll(expression, " ", "")
	tokens := calc.Tokenize(expression)
	tokens = calc.InfixToPostfix(tokens)
	if !calc.CountOp(tokens) {
		errors <- calc.ErrInvalidOperands
		done <- 1
		return
	}
	node := buildAST(tokens)
	result, err := calcLvl(node, tasks, results, id)
	if err != nil {
		errors <- err
	} else {
		results <- result
	}
	done <- 1
}
