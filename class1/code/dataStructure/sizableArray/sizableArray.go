package sizableArray

import "errors"

type SizeableArray struct {
	array []interface{}
	length int
	capacity int
}

// Init 初始化:容量固定 长度固定 底层数组所占用的内存空间大小也固定
func(s *SizeableArray) Init(length int) error {
	if length < 0 {
		return errors.New("illegal length")
	}

	s.length = length
	s.capacity = length
	s.array = make([]interface{}, length, length)
	return nil
}

// Insert 随机插入 定长数组初始化时已经将所有元素初始化为对应类型的零值
// 所以插入操作不影响长度和容量
func(s *SizeableArray) Insert(key int, val interface{}) error {
	if key >= s.length || key < 0 {
		return errors.New("illegal index")
	}

	s.array[key] = val
	return nil
}

// Delete 随机删除 go语言的定长数组中 删除一个元素实际上是把该元素赋值为对应类型的零值
// 因此对长度和容量没有影响
func(s *SizeableArray) Delete(key int) error {
	if key >= s.length || key < 0 {
		return errors.New("illegal index")
	}

	s.array[key] = nil
	return nil
}

// Lookup 随机访问
func (s *SizeableArray) Lookup(key int) (val interface{}, err error) {
	if key >= s.length || key < 0 {
		return nil, errors.New("illegal index")
	}

	return s.array[key], nil
}

func (s *SizeableArray) Length() int {
	return s.length
}

func (s *SizeableArray) Capacity() int {
	return s.capacity
}