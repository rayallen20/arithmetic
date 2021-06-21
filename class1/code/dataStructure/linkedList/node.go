package linkedList

type Node struct {
	value interface{}
	next *Node
}

func (n *Node) SetValue(value interface{})  {
	n.value = value
}

func (n Node) GetValue() interface{} {
	return n.value
}

func (n *Node) SetNext(next *Node) {
	n.next = next
}

func (n Node) GetNext() *Node {
	return n.next
}
