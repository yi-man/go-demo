package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
*/
func validParentheses(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	stack := []rune{}
	for _, char := range s {
		if char == '(' || char == '[' || char == '{' {
			stack = append(stack, char)
		} else {
			if len(stack) == 0 {
				return false
			} else if char == ')' && stack[len(stack)-1] == '(' || char == ']' && stack[len(stack)-1] == '[' || char == '}' && stack[len(stack)-1] == '{' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		}
	}
	return len(stack) == 0
}

func TestValidParentheses(t *testing.T) {

	result := validParentheses("{{[[(())]]}}")
	expected := true

	assert.Equal(t, result, expected, "TestValidParentheses test: {{[[(())]]}} true")
}
