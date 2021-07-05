package doubleLinkedList

import (
	"errors"
	"fmt"
)

type DoubleLinkedList struct {
	head *Node
	len  int
	tail *Node
}

func (d *DoubleLinkedList) Prepend(value interface{}) {
	node := Node{}
	node.SetValue(value)
	// 细节:DoubleLinkedList为空
	if d.head == nil {
		d.tail = &node
	} else {
		// 主体思路:
		// step1. 置Prepend前 DoubleLinkedList头部Node的Prev指针 指向待Prepend的节点
		// step2. 置待Prepend节点的Next指针为DoubleLinkedList的head指针
		// step3. 将DoubleLinkedList的head指针指向待Prepend的节点
		d.head.SetPrev(&node)
		node.SetNext(d.head)
	}

	d.head = &node
	d.len++
}

func (d *DoubleLinkedList) LookupFromHead() {
	node := d.head

	for i := 0; i <= d.len - 1; i++ {
		fmt.Printf("%v\n", node.GetValue())
		node = node.GetNext()
	}
}

func (d *DoubleLinkedList) LookUpFromTail() {
	node := d.tail
	for i := d.len - 1; i >= 0; i-- {
		fmt.Printf("%v\n", node.GetValue())
		node = node.GetPrev()
	}
}

func (d *DoubleLinkedList) Append(value interface{}) {
	node := Node{}
	node.SetValue(value)

	// 细节:DoubleLinkedList为空
	if d.head == nil {
		d.head = &node
	} else {
		// 主体思路:
		// step1. 置Append前 DoubleLinkedList尾部Node的Next指针 指向待Append的节点
		// step2. 置待Append节点的Prev指针为DoubleLinkedList的tail指针
		// step3. 将DoubleLinkedList的tail指针指向待Append的节点
		d.tail.SetNext(&node)
		node.SetPrev(d.tail)
	}

	d.tail = &node
	d.len++
}

func (d *DoubleLinkedList) Insert(index int, value interface{}) (err error) {

	// 细节:判定index是否合法
	if index < 0 || index > d.len {
		return errors.New("illegal index")
	}

	// 细节:index == 0 等价于 Prepend
	if index == 0 {
		d.Prepend(value)
		return nil
	}

	// 细节:index == d.len 等价于Append
	if index == d.len {
		d.Append(value)
		return nil
	}

	// 主体思路:
	target, err := d.find(index)
	if err != nil {
		return err
	}

	beInsertedNode := &Node{value: value}
	d.insert(beInsertedNode, target)

	return nil
}

func (d *DoubleLinkedList) find(index int) (target *Node, err error) {
	if index < 0 || index >= d.len {
		return nil, errors.New("illegal index")
	}

	// 细节:若index <= d.len / 2 则从前向后遍历 否则从后向前遍历
	if index <= d.len / 2 {
		target = d.head
		for i := 0; i < index; i++ {
			target = target.GetNext()
		}
	} else {
		target = d.tail
		for i := d.len - 1; i > index; i-- {
			target = target.GetPrev()
		}
	}

	return target, nil
}

func (d *DoubleLinkedList) insert(beInsertedNode, target *Node) {
	// step1. 寻找待插入索引处的节点(target)
	// step2. 待插入节点的Prev指针指向target.Prev
	// step3. target.Prev的Next指针指向待插入节点
	// step4. 待插入节点的Next指针指向target
	// step5. target的Prev指针指向待插入节点
	beInsertedNode.SetPrev(target.GetPrev())
	target.GetPrev().SetNext(beInsertedNode)
	beInsertedNode.SetNext(target)
	target.SetPrev(beInsertedNode)
	d.len++
}

func (d *DoubleLinkedList) Delete(index int) (err error) {
	// 细节:判定index是否合法
	if index < 0 || index >= d.len {
		return errors.New("illegal index")
	}

	// 细节:删除头部
	if index == 0 {
		d.deleteHead()
		return nil
	}

	// 细节:删除尾部
	if index == d.len - 1 {
		d.deleteTail()
		return nil
	}

	// 主体思路
	target, err := d.find(index)
	if err != nil {
		return err
	}

	d.deleteInternal(target)
	return nil
}

func (d *DoubleLinkedList) deleteHead() {
	d.head = d.head.GetNext()
	if d.head == nil {
		d.tail = nil
	} else {
		d.head.SetPrev(nil)
	}
	d.len--
}

func (d *DoubleLinkedList) deleteInternal(target *Node) {
	target.GetPrev().SetNext(target.GetNext())
	target.GetNext().SetPrev(target.GetPrev())
	d.len--
}

func (d *DoubleLinkedList) deleteTail() {
	d.tail = d.tail.GetPrev()
	if d.tail == nil {
		d.head = nil
	} else {
		d.tail.SetNext(nil)
	}
	d.len--
}