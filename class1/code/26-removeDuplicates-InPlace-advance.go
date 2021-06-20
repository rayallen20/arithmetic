package code

func removeDuplicates(nums []int) int {
	currentInsert := 0
	for i := 0; i <= len(nums) - 1; i++ {
		if i == 0 || nums[i] > nums[i - 1] {
			nums[currentInsert] = nums[i]
			currentInsert++
		}
		if i !=0 && i != len(nums) - 1 && nums[i] < nums[i - 1] {
			break
		}
	}
	return currentInsert
}