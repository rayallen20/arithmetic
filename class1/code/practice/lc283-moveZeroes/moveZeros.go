package lc283_moveZeroes

func MoveZeroes(nums []int)  {
	originLen := len(nums)
	zeroTarget := len(nums) - 1

	for i := 0; i < originLen; i++ {
		if nums[i] == 0 {
			copy(nums[i: originLen], nums[i + 1: originLen])
			nums[zeroTarget] = 0
			i--
			zeroTarget--
			originLen--
		}
	}
}

func MoveZeros2(nums []int) {
	nonzeroPosition := 0

	for i := 0; i < len(nums); i++ {
		// 元素值非0就要
		if nums[i] != 0 {
			nums[nonzeroPosition] = nums[i]
			nonzeroPosition++
		}
	}

	// 从nonzeroPosition处开始 直到数组结尾处 皆为要用0覆盖的位置
	for nonzeroPosition < len(nums) {
		nums[nonzeroPosition] = 0
		nonzeroPosition++
	}
}
