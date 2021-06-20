package code

import "fmt"

func merge(nums1 []int, m int, nums2 []int, n int) {
	result := make([]int, 0, m + n)

	// nums2是空slice
	if n == 0 {
		for _, v := range nums1 {
			result = append(result, v)
			fmt.Printf("%v\n", result)
		}
		return
	}

	// nums1是空slice
	if m == 0 {
		for _, v := range nums2 {
			result = append(result, v)
			fmt.Printf("%v\n", result)
		}
		return
	}

	j := 0

	for i := 0; i <= m - 1; i++ {
		if nums1[i] < nums2[j] {
			// i小j大 放i i+1 j不动
			result = append(result, nums1[i])
		} else {
			// j小i大 或 ij相等 放j j+1 i不动
			result = append(result, nums2[j])
			j++
			i--
		}

		// i循环完了还有j 把j都放到结果中
		if i == m - 1 && j <= n - 1 {
			for ;j <= n - 1; j++ {
				result = append(result, nums2[j])
			}
		}

		// j循环完了还有i 把i放到结果中
		if j > n - 1 {
			if i < 0 {
				for k, v := range nums1 {
					if k <= m - 1 {
						result = append(result, v)
					} else {
						break
					}
				}
			} else {
				for ; i < m - 1; i++ {
					result = append(result, nums1[i])
				}
			}
			break
		}
	}

	fmt.Printf("%v\n", result)
}