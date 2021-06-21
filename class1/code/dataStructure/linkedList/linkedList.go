package linkedList

import (
	"errors"
	"fmt"
)

type LinkedList struct {
	head *Node
	len  int
}

func (l *LinkedList) IsEmpty() bool {
	if l.head == nil {
		return true
	}
	return false
}

func (l *LinkedList) Prepend(value interface{}) {
	node := Node{}
	node.SetValue(value)
	if l.head == nil {
		l.head = &node
	} else {
		node.SetNext(l.head)
		l.head = &node
	}
	l.len++
}

func (l *LinkedList) Lookup() {
	node := l.head

	for i := 0; i <= l.len - 1; i++ {
		fmt.Printf("%v\n", node.GetValue())
		node = node.GetNext()
	}
}

func (l *LinkedList) Append(value interface{}) {
	tailNode := Node{}
	tailNode.SetValue(value)
	// 细节:若LinkedList为空 直接设置Head指针指向待Append的Node即可
	if l.len == 0 {
		l.head = &tailNode
	} else {
		// 主体思路:找到Append前的最后一个Node,将该Node的Next指针指向待Append的Node即可
		var lastNode *Node
		now := l.head

		for i := 0; i <= l.len - 1; i++ {
			if i == l.len - 1 {
				lastNode = now
			}
			now = now.GetNext()
		}
		lastNode.SetNext(&tailNode)
	}

	l.len++
}

func(l *LinkedList) Insert(index int, value interface{}) {
	// 细节:index == 0 等价于Prepend操作
	if index == 0 {
		l.Prepend(value)
	} else if index == l.len {
		// 细节:index == len 等价于Append操作
		l.Append(value)
	} else if index < 0 || index > l.len {
		// 细节:判定index是否合法
		err := errors.New("illegal index")
		fmt.Printf("%v\n", err)
	} else {
		// 主体思路:找到index位置前的Node 将待Insert的Node的Next指针设置为和该Node的Next指针相同
		// 再将该Node的Next指针指向待Insert的Node
		insertNode := Node{}
		insertNode.SetValue(value)

		prevNode := l.head

		for i := 0; i < index - 1; i++ {
			prevNode = prevNode.GetNext()
		}
		insertNode.SetNext(prevNode.GetNext())
		prevNode.SetNext(&insertNode)
		l.len++
	}
}

func (l *LinkedList) Delete(index int) {
	// 细节:判定l是否是一个空Linked List
	if l.head == nil {
		err := errors.New("can't delete element in a nil linked list")
		fmt.Printf("%v\n", err)
		return
	}

	// 细节:判定index的是否合法
	if index < 0 || index > l.len - 1 {
		err := errors.New("illegal index")
		fmt.Printf("%v\n", err)
	} else if index == 0 {
		// 细节: index == 0 置l.Head = l.Head.GetNext()即可
		l.head = l.head.GetNext()
		l.len--
	} else {
		// 主体思路:找到index位置前的Node 将该Node的Next指针指向待删除Node的Next指针
		prevNode := l.head
		for i := 0; i < index - 1; i++ {
			prevNode = prevNode.GetNext()
		}
		prevNode.SetNext(prevNode.GetNext().GetNext())
		l.len--
	}
}
