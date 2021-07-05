package doubleLinkedList

type Node struct {
	value interface{}
	next *Node
	prev *Node
}

func (n *Node) SetValue(value interface{}) {
	n.value = value
}

func (n *Node) GetValue() interface{} {
	return n.value
}

func (n *Node) SetNext(next *Node) {
	n.next = next
}

func (n *Node) GetNext() *Node {
	return n.next
}

func (n *Node) SetPrev(prev *Node) {
	n.prev = prev
}

func (n *Node) GetPrev() *Node {
	return n.prev
}
