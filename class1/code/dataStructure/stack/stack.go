package stack

import (
	"errors"
	"fmt"
)

type Stack struct {
	// 数组尾部为栈顶 头部为栈底
	items []interface{}
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (item interface{}, err error) {
	if s.IsEmpty() {
		return nil, errors.New("stack is empty")
	}

	item = s.items[len(s.items) - 1]
	s.items = s.items[:len(s.items) - 1]
	return item, nil
}

func (s *Stack) Peek() (item interface{}, err error) {
	if len(s.items) == 0 {
		return nil, errors.New("stack is empty")
	}

	item = s.items[len(s.items) - 1]
	return item, nil
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Check() {
	for i := len(s.items) - 1; i >= 0; i-- {
		fmt.Printf("%v\n", s.items[i])
	}
}
