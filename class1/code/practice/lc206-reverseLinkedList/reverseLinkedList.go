package lc206_reverseLinkedList

type ListNode struct {
	Val int
	Next *ListNode
}

func reverseLinkedList(head *ListNode) *ListNode {
	now := head
	var arr []*ListNode

	// 将链表中的所有节点按顺序存入数组
	for now != nil {
		arr = append(arr, now)
		now = now.Next
	}

	if len(arr) == 0 {
		return nil
	}

	// 设置数组中所有元素的Next指针指向前一个元素
	for i := len(arr) - 1; i >= 0; i-- {
		if i == 0 {
			arr[i].Next = nil
		} else {
			arr[i].Next = arr[i - 1]
		}
	}

	return arr[len(arr) - 1]
}

func reverseLinkedList2(head *ListNode) *ListNode {
	now := head
	var prev *ListNode

	for now != nil {
		// 暂存下一个
		next := now.Next

		// 置Next指针指向上一个
		now.Next = prev

		// 置自己为上一个
		prev = now

		// 处理下一个节点
		now = next
	}

	return prev
}
