package lc88_merge

func Merge(nums1 []int, m int, nums2 []int, n int)  {
	// 特例:nums1中无有效元素
	if m == 0 {
		for k, v := range nums2 {
			nums1[k] = v
		}
		return
	}

	// nums2中剩余未操作的元素标记
	var residueNums2 int

	// nums1中上一次发生位移的元素标记
	var lastMoveNums1 int

	// 对于nums2 从上一次未操作的元素开始遍历 查找是否在nums1中有该元素合适的位置
	for i := residueNums2; i < n; i++ {
		for j := lastMoveNums1; j < m; j++ {
			if nums1[j] > nums2[i] {
				moveBack(nums1, j, m - 1)
				nums1[j] = nums2[i]
				lastMoveNums1 = j
				m++
				residueNums2++
				break
			}
		}

		// 若此时nums2的元素 比nums1中最大的元素值还要大 说明此时nums2中剩余的元素都不需要再在nums1中寻找插入位置了
		if nums2[i] > nums1[m - 1] {
			break
		}
	}

	// 把nums2中剩余未被处理的元素放到nums1的末尾
	if residueNums2 < n {
		for k := residueNums2; k < n; k++ {
			nums1[m] = nums2[k]
			m++
		}
	}
}

func moveBack(nums []int, start, end int) {
	for i := end + 1; i > start; i-- {
		nums[i] = nums[i - 1]
	}
}
