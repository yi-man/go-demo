package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func twoSum(nums []int, target int) []int {
	hashMap := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		if _, ok := hashMap[target-nums[i]]; ok {
			return []int{hashMap[target-nums[i]], i}
		} else {
			hashMap[nums[i]] = i
		}
	}
	return []int{-1, -1}
}

func TestTwoSum(t *testing.T) {
	result := twoSum([]int{2, 7, 11, 15}, 9)
	expected := []int{0, 1}

	assert.Equal(t, result, expected, "twoSum([]int{2, 7, 11, 15}, 9) should be equal to []int{0, 1}")
}
