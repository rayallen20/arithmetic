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
