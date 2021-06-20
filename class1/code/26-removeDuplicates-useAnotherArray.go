package code

func removeDuplicates(nums []int) []int {
	result := make([]int, 0, len(nums) - 1)
	for i := 0; i <= len(nums) - 1; i++ {
		if i == 0 || nums[i] != nums[i - 1] {
			result = append(result, nums[i])
		}
	}
	return result
}