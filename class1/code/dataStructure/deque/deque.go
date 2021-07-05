package deque

import (
	"errors"
	"fmt"
)

type Deque struct {
	// 数组头部表示双端队列头部 数组尾部表示双端队列尾部
	deque []interface{}
}

func (d *Deque) IsEmpty() bool {
	return len(d.deque) == 0
}

func(d *Deque) AddRear(item interface{}) {
	d.deque = append(d.deque, item)
}

func (d *Deque) AddFront(item interface{}) {
	tmp := []interface{}{item}
	tmp = append(tmp, d.deque...)
	d.deque = tmp
}

func (d *Deque) Size() int {
	return len(d.deque)
}

func (d *Deque) RemoveRear() (item interface{}, err error) {
	if d.IsEmpty() {
		return nil, errors.New("deque is empty")
	}

	item = d.deque[len(d.deque) - 1]
	d.deque = d.deque[0:len(d.deque) - 1]
	return item, nil
}

func (d *Deque) RemoveFront() (item interface{}, err error) {
	if d.IsEmpty() {
		return nil, errors.New("deque is empty")
	}

	item = d.deque[0]
	d.deque = d.deque[1:]
	return item, nil
}

func (d *Deque) Check() {
	for i := 0; i < len(d.deque); i++ {
		fmt.Printf("%v\n", d.deque[i])
	}
}
