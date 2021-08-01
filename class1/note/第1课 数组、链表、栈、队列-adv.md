# 第1课 数组、链表、栈、队列

## 1. 数组

数组在内存中是一段连续的存储空间

数组的基本特点:**支持随机访问**

数组的关键操作:索引与寻址

### 1.1 定长数组

#### 1.1.1 定长数组的定义

初始化时给定一个整数N,该整数表示数组最大容量,数组最多只能容纳N个元素.定长数组不具备自动扩缩容的能力

#### 1.1.2 定长数组的实现

```go
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
```

### 1.2 变长数组

#### 1.2.1 变长数组的定义

定长数组的长度不可改变,这个特性使得在特定的场景中,这样的集合不太适用.因此需要一种更为通用的"动态数组".与定长数组相比,变长数组的长度是不固定的,可以追加元素,在追加时,可能使变长数组的容量增大.

#### 1.2.2 变长数组的实现

在实现变长数组时,就会有以下问题要解决:

- 需要支持索引与随机访问
- 分配多长的连续空间合适?
- 空间不够用了怎么办?
- 空间剩余很多如何回收?

此处我们可以参考GO语言中slice的实现.

GO语言对slice的定义:`$GOROOT/src/runtime/slice.go`

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

GO语言中slice的扩容机制:`$GOROOT/src/runtime/slice.go`

```go
	if old.len < 1024 {
		newcap = doublecap
	} else {
		// Check 0 < newcap to detect overflow
		// and prevent an infinite loop.
		for 0 < newcap && newcap < cap {
			newcap += newcap / 4
		}
		// Set newcap to the requested cap when
		// the newcap calculation overflowed.
		if newcap <= 0 {
			newcap = cap
		}
	}
```

但如果我们自己去实现,起手面临的问题在于:如果使用底层数组作为变长数组的基础,那么GO语言中,长度是数组的一部分.无法实现扩缩容.因此我们变长数组的实现,使用切片代替底层数组.通过变长数组的长度和容量实现对这个切片的访问控制.

```go
package resizableArray

import (
	"errors"
	"fmt"
)

// growFactor 扩容因子 每次扩容后的容量大小为原来的2倍
const growFactor = 2

// reduceFactor 缩容因子 每次缩容后的容量大小为原来的1/2
const reduceFactor = 0.5

// reduceThreshold 缩容阈值 当 ResizableArray.len / ResizableArray.cap <= reduceThreshold 时 触发缩容机制
const reduceThreshold = 0.25

// ResizableArray 变长数组
type ResizableArray struct {
	// array 最规范的做法是使用定长数组作为变长数组的底层 但定长数组的长度一旦定义就无法改变
	// array 作为变长数组结构体中的一个字段 在扩缩容时 需要改变长度 所以使用了slice作为该字段的类型
	// len(array) 表示 ResizableArray 的长度
	array []interface{}

	// len 变长数组的长度
	len int

	// cap 变长数组的容量 实际上就是 len(array)
	cap int
}

// Init 初始化时给定的容量将作为底层数组的长度使用
// 而底层数组的长度又是变长数组的容量
func (r *ResizableArray) Init(cap int) error {
	if cap < 0 {
		return errors.New("illegal cap")
	}
	r.array = make([]interface{}, cap, cap)
	r.cap = cap
	return nil
}

// Prepend 头部插入
func (r *ResizableArray) Prepend(val interface{}) {
	r.grow()
	r.moveBack(0, r.len - 1, 1)
	r.array[0] = val
	r.len++
}

// Append 后向追加
func (r *ResizableArray) Append(val interface{}) {
	r.grow()

	r.array[r.len] = val
	r.len++
}

func (r *ResizableArray) Insert(index int, val interface{}) error {
	if 0 > index || index > r.len {
		return errors.New("illegal index")
	}

	r.grow()
	r.moveBack(index, r.len - 1, 1)
	r.array[index] = val
	r.len++

	return nil
}

func (r *ResizableArray) Delete(index int) (val interface{}, err error) {
	if index < 0 || index > r.len - 1 {
		return nil, errors.New("illegal index")
	}

	// 将要删除的元素暂存 以便后续返回
	res := make([]interface{},1, 1)
	copy(res, r.array[index:index+1])
	r.moveFront(index + 1, r.len, 1)
	r.len--

	r.reduce()
	return res[0], nil
}

// grow 检测是否需要扩容 如果需要 则把原底层数组修改为扩容后的底层数组 并维护变长数组的cap字段
// 扩容机制:若原底层数组长度为0 则扩容后底层数组长度为1 若原底层数组长度不为0 则扩容后底层数组长度为原来的2倍
func (r *ResizableArray) grow() {
	afterAddLen := r.len + 1
	if afterAddLen <= r.cap {
		return
	}

	// 添加元素会触发扩容机制
	var newArray []interface{}
	var newCap int

	// 创建新的底层数组
	if r.cap == 0 {
		newCap = 1
	} else {
		newCap = r.cap * growFactor
	}

	newArray = make([]interface{}, newCap, newCap)
	// 将原底层数组的元素拷贝到新的底层数组中
	copy(newArray, r.array)

	r.cap = newCap
	r.array = newArray
}

// reduce 检测是否需要缩容 如果需要 则把原底层狐族修改为缩容后的底层数组 并维护变长数组的cap字段
// 缩容机制触发条件: 当 ResizableArray.len / ResizableArray.cap <= 25%时 触发缩容机制
// 缩容机制:释放一半空间 即新的底层数组的长度为原底层数组长度的一半
func (r *ResizableArray) reduce()  {
	nowThreshold := float64(r.len) / float64(r.cap)

	if nowThreshold <= reduceThreshold {
		r.cap = int(float64(r.cap) * reduceFactor)
		newArray := make([]interface{}, r.cap, r.cap)
		copy(newArray, r.array)
		r.array = newArray
	}
}

// moveBack 从给定的start位置开始 到给定的end位置为止 向后移动offset单位长度
func (r *ResizableArray) moveBack(start, end, offset int) {
	for i := end + 1; i > start; i-- {
		r.array[i] = r.array[i - offset]
	}
}

// moveFront 从给定的start位置开始 到给定的end位置为止 向前移动offset单位长度
func (r ResizableArray) moveFront(start, end, offset int)  {
	for i := start; i < end; i++ {
		r.array[i - offset] = r.array[i]
	}
}

// Lookup 随机访问
func (r *ResizableArray) Lookup(index int) (val interface{}, err error) {
	if index < 0 || index > r.len - 1 {
		return nil, errors.New("illegal index")
	}

	return r.array[index], nil
}

func (r ResizableArray) Check()  {
	fmt.Printf("%v\n", r.array[:r.len])
	fmt.Printf("resizableArray Len:%d resizableArray Cap:%d\n", r.len, r.cap)
}
```



















