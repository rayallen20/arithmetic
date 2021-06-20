package code

func merge(nums1 []int, m int, nums2 []int, n int) {
	i := m - 1
	j := n - 1
	checkIndex := m + n - 1

	for ; checkIndex >= 0; checkIndex-- {
		// 主体思路:对nums1和nums2 从后向前遍历 二者中较大的值 放到nums1的尾部
		if j < 0 || (i >= 0 && nums1[i] > nums2[j]) {
			nums1[checkIndex] = nums1[i]
			i--
		} else {
			// 等价于 j >= 0 && !(i >= 0 && nums1[i] > nums2[j])
			// 等价于 j >= 0 && (i < 0 || nums1[i] <= nums2[j])
			// 本段逻辑同样适用于当i < 0 且nums2中还有没处理的元素时(即j > 0时)
			nums1[checkIndex] = nums2[j]
			j--
		}
	}
}