package main

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// nums 不唯一，有可能相等
func threeNumSum(nums []int, target int) [][]int {
	result := [][]int{}

	sort.Ints(nums)

	if len(nums) < 3 {
		return result
	}

	// 遍历
	for i := 0; i < len(nums)-2; i++ {
		base := nums[i]

		// base大于目标时，中断循环
		if base > target {
			break
		}
		// 过滤掉相等的值
		if i > 0 && base == nums[i-1] {
			continue
		}

		innerTarget := target - base

		// 从i + 1 开始双指针遍历
		low := i + 1
		high := len(nums) - 1

		for low < high {
			// 相等时，
			if nums[low]+nums[high] == innerTarget {
				record := []int{base, nums[low], nums[high]}
				result = append(result, record)

				// 去重
				for low < high && nums[low] == nums[low+1] {
					low++
				}
				for low < high && nums[high] == nums[high-1] {
					high--
				}

				// 继续看有没有满足的值
				low++
				high--
				// break
			} else if nums[low]+nums[high] < innerTarget {
				low++
			} else {
				high--
			}
		}

	}

	return result
}

func TestThreeNumSum(t *testing.T) {
	result := threeNumSum([]int{-100, -10, 200, 100, 0, 2, -2, 19, -10, -9, 4, 6}, 0)
	expected := [][]int{{-100, 0, 100}, {-10, -9, 19}, {-10, 4, 6}, {-2, 0, 2}}

	assert.Equal(t, result, expected, "threeNumSum test")
}
