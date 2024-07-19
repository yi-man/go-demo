package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode"
)

// applyOperator 从操作符栈中弹出一个操作符，并从操作数栈中弹出两个操作数进行计算，将结果重新压入操作数栈。
func applyOperator(operators *[]rune, values *[]int) {
	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	right := (*values)[len(*values)-1]
	*values = (*values)[:len(*values)-1]
	left := (*values)[len(*values)-1]
	*values = (*values)[:len(*values)-1]

	switch operator {
	case '+':
		*values = append(*values, left+right)
	case '-':
		*values = append(*values, left-right)
	case '*':
		*values = append(*values, left*right)
	case '/':
		*values = append(*values, left/right)
	}
}

// precedence 返回操作符的优先级。
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

// evaluateExpression 处理包含括号的数学表达式，返回计算结果。
func evaluateExpression(expression string) int {
	nums := []int{}
	ops := []rune{}

	i := 0
	for i < len(expression) {
		char := rune(expression[i])
		if char == ' ' {
			i++
			continue
		}

		if char == '(' {
			ops = append(ops, char)
		} else if unicode.IsDigit(char) {
			val := 0
			for i < len(expression) && unicode.IsDigit(rune(expression[i])) {
				val = val*10 + int(expression[i]-'0')
				i++
			}
			nums = append(nums, val)
			i-- // 因为for循环会增加i
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				applyOperator(&ops, &nums)
			}
			ops = ops[:len(ops)-1] // 弹出 '('
		} else {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(char) {
				applyOperator(&ops, &nums)
			}
			ops = append(ops, char)
		}
		i++
	}

	for len(ops) > 0 {
		applyOperator(&ops, &nums)
	}

	return nums[0]
}

func TestExpression(t *testing.T) {

	result := evaluateExpression("3 + (2 * 2) - 1")
	expected := 6

	assert.Equal(t, result, expected, "TestExpression test: 3 + (2 * 2) - 1=6")
}
