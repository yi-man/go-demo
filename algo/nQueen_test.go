package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func nQueen(n int) [][]string {
	var result [][]string
	var board = make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var cols = make([]bool, n)
	var diag1 = make([]bool, 2*n-1)
	var diag2 = make([]bool, 2*n-1)

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			var solution []string
			for _, b := range board {
				solution = append(solution, string(b))
			}
			result = append(result, solution)
			return
		}
		for col := 0; col < n; col++ {
			if cols[col] || diag1[row+col] || diag2[row-col+n-1] {
				continue
			}
			board[row][col] = 'Q'
			cols[col] = true
			diag1[row+col] = true
			diag2[row-col+n-1] = true
			backtrack(row + 1)
			board[row][col] = '.'
			cols[col] = false
			diag1[row+col] = false
			diag2[row-col+n-1] = false
		}
	}

	backtrack(0)
	return result
}

func TestNQueen(t *testing.T) {

	result1 := nQueen(4)
	expected1 := [][]string{
		{".Q..", "...Q", "Q...", "..Q."},
		{"..Q.", "Q...", "...Q", ".Q.."},
	}
	assert.ElementsMatch(t, result1, expected1, "nQueen 4")

}
