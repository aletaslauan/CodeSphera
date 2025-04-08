package main

import "fmt"

func TwoSum(nums []int, target int) []int {
    // Create a map to store the numbers and their indices
    numMap := make(map[int]int)

    for i, num := range nums {
        // Calculate the complement of the current number
        complement := target - num

        // Check if the complement exists in the map
        if idx, exists := numMap[complement]; exists {
            return []int{idx, i}
        }

        numMap[num] = i
    }
    return nil
}

func main() {
    //nums := []int{2, 7, 11, 15}
    //target := 9
    //result := twoSum(nums, target)
    fmt.Println("TwoSum") 
}

/*
Two Sum Problem

Write a function that takes a slice of integers and a target integer,
and returns the indices of the two numbers that add up to the target.

Constraints:
- Do not use the same element twice
- Exactly one solution is guaranteed
- You can return the result in any order	

Example:
nums := []int{2, 7, 11, 15}
target := 9
Output: [0, 1] because nums[0] + nums[1] == 9
*/

// TwoSum finds two indices in nums whose values add up to target.
// It returns the indices as a slice of two integers.
// Assumes exactly one solution exists.
