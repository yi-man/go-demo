package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func swap(s []rune, i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func quanPailie(s []rune) []string {
	var ret []string

	var backtrace func(start int)
	backtrace = func(start int) {
		if start == len(s) {
			ret = append(ret, string(s))
			return
		}
		for i := start; i < len(s); i++ {
			swap(s, start, i)
			backtrace(start + 1)
			swap(s, start, i)
		}
	}

	backtrace(0)

	return ret
}

func permuteQuanPailie(chars []rune) []string {
	var result []string
	var temp []rune
	var used = make([]bool, len(chars))

	var backtrack func([]rune, []rune, []bool)
	backtrack = func(chars []rune, temp []rune, used []bool) {
		if len(temp) == len(chars) {
			result = append(result, string(temp))
			return
		}
		for i := 0; i < len(chars); i++ {
			if used[i] {
				continue
			}
			used[i] = true
			temp = append(temp, chars[i])
			backtrack(chars, temp, used)
			temp = temp[:len(temp)-1]
			used[i] = false
		}
	}

	backtrack(chars, temp, used)
	return result
}

func TestQuanPaiLie(t *testing.T) {
	input := "abc"
	chars := []rune(input)

	result1 := quanPailie(chars)
	expected1 := []string{"acb", "abc", "bac", "bca", "cab", "cba"}

	assert.ElementsMatch(t, result1, expected1, "quanPailie([]string{\"a\", \"b\", \"c\"})")

	result := permuteQuanPailie(chars)
	expected := []string{"abc", "acb", "bac", "bca", "cab", "cba"}
	assert.Equal(t, result, expected, "permuteQuanPailie([]string{\"a\", \"b\", \"c\"})")
}
