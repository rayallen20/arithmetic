package code

import "fmt"

func moveZeroes(nums []int)  {
	result := make([]int, 0, len(nums) - 1)

	for _, v := range nums {
		if v != 0 {
			result = append(result, v)
		}
	}

	for resultLen := len(result) - 1;resultLen < len(nums) - 1; resultLen++ {
		result = append(result, 0)
	}
	fmt.Printf("%v\n", result)
}