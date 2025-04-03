package main

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected []int // order doesn't matter
	}{
		{
			name:     "Basic example",
			nums:     []int{2, 7, 11, 15},
			target:   9,
			expected: []int{0, 1},
		},
		{
			name:     "Unordered solution",
			nums:     []int{3, 2, 4},
			target:   6,
			expected: []int{1, 2},
		},
		{
			name:     "Negative numbers",
			nums:     []int{-3, 4, 3, 90},
			target:   0,
			expected: []int{0, 2},
		},
		{
			name:     "Duplicates in input",
			nums:     []int{3, 3},
			target:   6,
			expected: []int{0, 1},
		},
		{
			name:     "Target is sum of first and last",
			nums:     []int{1, 2, 3, 4, 5, 6},
			target:   7,
			expected: []int{0, 5},
		},
		{
			name:     "Middle elements add up",
			nums:     []int{1, 4, 5, 6},
			target:   9,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TwoSum(tt.nums, tt.target)

			// Check if both elements are in result, regardless of order
			if len(result) != 2 {
				t.Fatalf("Expected 2 indices, got %d", len(result))
			}

			i, j := result[0], result[1]
			if tt.nums[i]+tt.nums[j] != tt.target {
				t.Errorf("nums[%d] + nums[%d] = %d, want %d", i, j, tt.nums[i]+tt.nums[j], tt.target)
			}

			// Optional strict order check:
			if !reflect.DeepEqual(result, tt.expected) && !reflect.DeepEqual([]int{result[1], result[0]}, tt.expected) {
				t.Logf("Returned indices: %v (order may vary)", result)
			}
		})
	}
}

