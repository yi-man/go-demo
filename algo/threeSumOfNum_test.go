package main

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func threeNumSum(nums []int, target int) [][]int {
	result := [][]int{}

	sort.Ints(nums)

	if len(nums) < 3 {
		return result
	}

	return result
}

func TestThreeNumSum(t *testing.T) {
	result := threeNumSum([]int{-1, 0, 1, 2, -1, -4}, 1)
	expected := [][]int{{-1, -1, 2}, {-1, 0, 1}}

	assert.Equal(t, result, expected, "threeNumSum([]int{-1, 0, 1, 2, -1, -4}, 1) should be equal to [][]int{{-1, -1, 2}, {-1, 0, 1}}")
}
