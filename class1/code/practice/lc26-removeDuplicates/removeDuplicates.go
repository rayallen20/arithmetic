package lc26_removeDuplicates

func RemoveDuplicates(nums []int) int {
	// 空数组或只有1个元素的数组 必然没有重复元素
	if len(nums) == 0 || len(nums) == 1 {
		return len(nums)
	}

	// 数组中至少有2个元素 此时第1个元素的值必然为非重复值
	var elementsCounter int = 1

	for i := 1; i < len(nums); i++ {
		// 后一个元素的值 > 前一个元素的值 则后一个元素必然是非重复元素
		if nums[i] > nums[i - 1] {
			elementsCounter++
		} else {
			// 后一个元素的值 <= 前一个元素的值 则后一个元素是重复元素
			// 从该元素向后 找到第1个 > 前一个元素值的元素 覆盖
			found := false
			for j := i + 1; j < len(nums); j++ {
				if nums[j] > nums[i - 1] {
					nums[i] = nums[j]
					elementsCounter++
					found = true
					break
				}
			}

			// 若向后没有找到 > 前一个元素值的元素 则说明后面所有元素都比前一个元素小
			// 说明后面都是重复元素了
			if !found {
				return elementsCounter
			}
		}
	}

	return elementsCounter
}
