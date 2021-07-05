package queue

import (
	"errors"
	"fmt"
)

type Queue struct {
	// 数组尾部为队头 头部为队尾(个人感觉这样比尾部表示队尾 头部表示队头好写)
	queue []interface{}
}

func (q *Queue) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *Queue) Enqueue(item interface{}) {
	q.queue = append(q.queue, item)
}

func (q *Queue) Dequeue() (item interface{}, err error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	item = q.queue[0]
	q.queue = q.queue[1:]
	return item, nil
}

func (q *Queue) Size() int {
	return len(q.queue)
}

func (q *Queue) Check() {
	for i := len(q.queue) - 1; i >= 0; i-- {
		fmt.Printf("%v\n", q.queue[i])
	}
}
