package code

func moveZeroes(nums []int) {
	// 找到数组中第1个值为0的元素
	currentIndex := 0
	for ;currentIndex <= len(nums) - 1; currentIndex++ {
		if nums[currentIndex] == 0 {
			break
		}
	}


	// 如果没找到 说明数组不需要修改
	if currentIndex - 1 == len(nums) - 1 {
		return
	}
	
	zeroCounter := 1

	for i := currentIndex + 1; i <= len(nums) - 1; i++ {
		if nums[i] != 0 {
			nums[currentIndex] = nums[i]
			currentIndex++
		} else {
			zeroCounter++
		}
	}

	// 用0覆盖结尾处应该为0的元素
	for ;zeroCounter > 0; zeroCounter-- {
		nums[len(nums) - zeroCounter] = 0
	}
}
