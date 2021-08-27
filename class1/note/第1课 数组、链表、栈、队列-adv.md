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

// TODO:我的设计 为什么用25%作为trigger 50%作为reduceFactor

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

### 1.3 相关习题

#### 1.3.1 合并两个有序数组

[合并两个有序数组](https://leetcode-cn.com/problems/merge-sorted-array/)

##### a. 题目要求

给你2个有序整数数组`nums1`和`nums2`,请你将`nums2`合并到`nums1`中,使`nums1`成为一个有序数组.

初始化`nums1`和`nums2`的元素数量分别为`m`和`n`.你可以假设`nums1`的空间大小等于`m + n`,这样它就有足够的空间保存来自`nums2`的元素.

示例1:

```

输入: nums1 = [1, 2, 3, 0, 0, 0], m = 3, nums2 = [2, 5, 6], n = 3

输出: [1, 2, 3, 4, 5, 6]

```

示例2:

```

输入: nums1 = [1], m = 1, nums2 = [], n = 0

输出: [1]

```

提示:

- nums1.length == m + n
- nums2.length = n
- 0 <= m, n <= 200
- 1 <= m + n <= 200
- -10^9 <= nums1[i], nums2[i] <= 10^9


```go

func merge(nums1 []int, m int, nums2 []int, n int)  {

}

```

##### b. 审题

抽象一点这道题用一句话形容:把`nums2`中的元素全部放入`nums1`中,并保持升序.

那么对于`nums2`中的元素,就有2种可能性:

1. 需要放到`nums1`的中间.即放入`nums1[0]`到`nums1[m]`这个区间内.这种元素需要满足条件:`nums1[0] <= nums2[n] <= nums1[m]`
2. 需要放到`nums1`的结尾.即放入`nums1[m + 1]`到`nums1[len(nums1) - 1]`这个区间内.这种元素需要满足条件:`nums1[m] < nums2[n]`

分析完了`nums2`,再回过头来看`nums1`.对于`nums1`来讲,参数`m`表示的是`nums1`中有效元素的数量.所谓有效元素,即`nums1`中需要参与排序的元素的个数.按照题意,`nums1.length == m + n`.说明`nums1`中的元素,实际上也有2种可能性:

1. 需要参与排序的.即`nums1[0]`到`nums1[m]`这个区间内的元素
2. 不需要参与排序的.`nums1[m + 1]`到`nums1[len(nums1) - 1]`这个区间内的元素.即这部分元素实际上起到了一个占位的作用.它们存在的意义是为了保证`nums1`有足够的长度容纳`nums2`的元素

但想要实现这个需求,至少是需要把`nums1`的有效元素部分和`nums2`都扫一遍的

##### c. 根据审题结果,寻找合适的数据结构

这道题实际上考察的并不是对数据结构的应用,而是对算法的设计.很明显这道题要用双指针.问题在于怎么用.

##### d. 实现思路

本题实际上的含义是:对于`nums2`中的每一个元素,在`nums1`中找到一个能够保持升序的位置,并将nums2中的元素放到这个位置上.由于给定的2个数组都是已经按升序排好序的,所以我们按顺序,先处理`nums2`中 <= `nums1[m - 1]`的元素,再处理`nums2`中 > `nums1[m - 1]`的元素

step1. 处理`nums2`中 <= `nums1[m - 1]`的元素

- 对于每个`nums2`中的元素,在`nums1`中寻找第1个大于该元素值的索引
	- 若找到,则把`nums1`中从该索引处开始,直到`m - 1`处为止的所有元素,向后挪动1个单位
	- 用`nums2`中当前的元素替代`nums1`中当前指针指向的元素
	- 对于`nums1`而言,此时有效元素增加了1个,所以表示有效元素个数的`m`要做`+1`的操作
	- 对于`nums2`而言,因为后续可能出现 > `nums1[m - 1]`的元素,所以需要记录一个索引.该索引之前的所有元素,均已经被插入到`nums1`中,该索引及其之后的所有元素,均未放到`nums1`中

step2. 处理`nums2`中 > `nums1[m - 1]`的元素

- 由于之前记录了`nums1`中的有效元素个数,也记录了`nums2`中尚未处理的元素位置.所以这一步只需要将`nums2`中所有尚未处理的元素,按顺序放入`nums1`的有效元素后面即可.

##### e. 实现

```go
package lc88_merge

func Merge(nums1 []int, m int, nums2 []int, n int)  {
	var residue int
	
	// nums2中需要插入到nums1中的元素
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if nums1[j] > nums2[i] {
				moveBack(nums1, j, m - 1)
				nums1[j] = nums2[i]
				residue++
				m++
				break
			}
		}
	}

	// nums2中还有没放到nums1中的元素
	if residue < n {
		for k := residue; k < n; k++ {
			nums1[m] = nums2[k]
			m++
		}
	}
}

func moveBack(nums []int, start, end int) {
	for i := end + 1; i > start; i-- {
		nums[i] = nums[i - 1]
	}
}
```

##### f. 改进

可以改进的地方有2个点:

1. 对于`nums1`来讲,不需要每一次都从0开始扫,因为给定的2个数组都是升序数组,所以应该从上一次发生插入的位置开始本次遍历
2. 对于`nums2`来讲,不需要遍历到末尾.只需要遍历到第1个 > `nums[m]`的元素,就可以终止遍历了

改进后的实现:

```go
package lc88_merge

func Merge(nums1 []int, m int, nums2 []int, n int)  {
	// 特例:nums1中无有效元素
	if m == 0 {
		for k, v := range nums2 {
			nums1[k] = v
		}
		return
	}

	// nums2中剩余未操作的元素标记
	var residueNums2 int

	// nums1中上一次发生位移的元素标记
	var lastMoveNums1 int

	// 对于nums2 从上一次未操作的元素开始遍历 查找是否在nums1中有该元素合适的位置
	for i := residueNums2; i < n; i++ {
		for j := lastMoveNums1; j < m; j++ {
			if nums1[j] > nums2[i] {
				moveBack(nums1, j, m - 1)
				nums1[j] = nums2[i]
				lastMoveNums1 = j
				m++
				residueNums2++
				break
			}
		}

		// 若此时nums2的元素 比nums1中最大的元素值还要大 说明此时nums2中剩余的元素都不需要再在nums1中寻找插入位置了
		if nums2[i] > nums1[m - 1] {
			break
		}
	}

	// 把nums2中剩余未被处理的元素放到nums1的末尾
	if residueNums2 < n {
		for k := residueNums2; k < n; k++ {
			nums1[m] = nums2[k]
			m++
		}
	}
}

func moveBack(nums []int, start, end int) {
	for i := end + 1; i > start; i-- {
		nums[i] = nums[i - 1]
	}
}
```

时间复杂度:

- 对于`nums2`中的元素,要么插入到`nums1`中,要么赋值到`nums1`中,所以插入+赋值这两个操作的时间复杂度为`n`.
- 对于`nums2`中所有要插入`nums1`的元素,每做一次插入操作,都需要`nums1`中大于该元素值的所有元素向后位移1个单位.设`nums2`中共有`x`个需要插入到`nums1`中的元素,每个元素插入至`nums1`时,向后位移操作需要进行`y`次.则位移操作的时间复杂度为`x * y`.

所以再次改进的重点,应该是去掉向后位移的操作.这部分操作占有大量的时间复杂度.

##### g. 改进-去掉向后位移的操作

要想没有向后位移的操作,就需要新开辟一个数组,该数组容量为`m + n`.遍历`nums1`和`nums2`,将二者中值较小的元素放入这个新开辟的数组中即可.

```go
package lc88_merge

func merge(nums1 []int, m int, nums2 []int, n int) {
	res := make([]int, m + n, m + n)
	nums1Now := 0
	nums2Now := 0

	for i := 0; i < m + n; i++ {
		// 放nums1中的元素 有2种可能性
		// 1. nums1[nums1Now] <= nums2[nums2Now]
		// 2. nums2Now == n 且 nums1Now < m 即nums2中的元素都放完了且nums1中还有没放的元素 后面只放nums1中的元素了
		if nums1Now < m && (nums2Now == n || nums1[nums1Now] <= nums2[nums2Now]) {
			// fmt.Printf("n2Now = %d\n", nums2Now)
			res[i] = nums1[nums1Now]
			nums1Now++
		} else {
			// fmt.Printf("n2Now = %d\n", nums2Now)
			res[i] = nums2[nums2Now]
			nums2Now++
		}
	}

	copy(nums1, res)
}
```

时间复杂度:

- 遍历`nums1`和`nums2`的时间复杂度为`m + n`
- 赋值操作的时间复杂度为`m + n`
- 故时间复杂度为O(2(m + n)),即O(m + n)

但此方法由于新开辟了一个长度和容量均为`m + n`的数组,所以空间复杂度较高

##### h. 改进-不需向后位移且不需开辟新数组的方案

我们来回顾f小节的算法.该算法之所以需要向后位移,是为了保留`nums1`中还未处理的元素,不让比较结果直接覆盖输入.若想要不做向后位移的操作,就只能是**在结果覆盖输入之前,让输入先于结果参与比较.**放在本题中即为:

- 从`nums1`的末尾处开始放入元素
- 将`nums1`和`nums2`中较大的元素放入

```go
package lc88_merge

func merge(nums1 []int, m int, nums2 []int, n int) {
	// 从nums1的末尾处开始放 二者中谁大放谁
	nums1Now := m - 1
	nums2Now := n - 1

	for i := m + n - 1; i >= 0; i-- {
		// 当前位置放放nums1中的元素 有2种情况:
		// 1. nums1[nums1Now] > nums2[nums2Now]
		// 2. nums2中的元素放完了且nums1中还有没放的元素
		if nums2Now < 0 || (nums1Now >= 0 && nums1[nums1Now] > nums2[nums2Now]) {
			nums1[i] = nums1[nums1Now]
			nums1Now--
		} else {
			nums1[i] = nums2[nums2Now]
			nums2Now--
		}
	}
}
```

#### 1.3.2 删除有序数组中的重复项

[删除有序数组中的重复项](https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array/)

##### a. 题目要求

给你一个有序数组`nums`,请你**原地**删除重复出现的元素,使每个元素**只出现一次**,返回删除后数组的新长度.

不要使用额外的数组空间,你必须在**原地**修改输入数组并在使用O(1)额外空间的条件下完成.

示例1:

```

输入: nums = [1, 1, 2]

输出: 2, nums = [1, 2]

解释: 函数应该返回新的长度为2,并且原数组nums的前2个元素被修改为1, 2.不需要考虑数组中超出新长度后面的元素

```

示例2:

```

输入: nums = [0, 0, 1, 1, 1, 2, 2, 3, 3, 4]

输出: 5, nums = [0, 1, 2, 3, 4]

解释: 函数应该返回新的长度为5,并且原数组nums的前5个元素被修改为0, 1, 2, 3, 4.不需要考虑数组中超出新长度后面的元素

```

提示:

-  0 <= nums.length <= 3 * 10^4
-  -10^4 <= nums[i] <= 10^4
-  nums已按升序排列

##### b. 审题

本题要求共有2点:

1. 返回数组中非重复元素的个数n
2. 数组中前n项为非重复项,并保持升序

本题给定的初始条件中,比较重要的信息有2点:

1. **给定数组是升序排列的**
2. **不需要考虑数组中超出新长度后面的元素**

因此,我们定义本场景下的基础操作为:

1. 找到重复元素
2. 找到该重复元素之后的第1个非重复元素
3. 用这个非重复元素**覆盖**重复元素

##### c. 根据审题结果,寻找合适的数据结构

从审题的分析结果中可知,需要2个指针:一个指向重复元素,另一个指向非重复元素.即**双指针**

##### d. 实现思路

- step1. 排除异常情况
	- 异常情况:给定数组长度为0或为1时,此时数组中必然没有重复元素,也不需要做去重.此时直接返回数组长度即可

- step2. 确认重复元素的索引
	- 由于题目中给出的条件为**给定数组是升序排列的**,所以**当前元素的值 <= 当前元素的前一个元素值**时,当前元素即为重复元素.

- step3. 确认重复元素之后的第1个非重复元素索引
	- 同样地,由于题目中给出的条件为**给定数组是升序排列的**,所以**从上一步找到的重复元素的索引向后,第1个 > 重复元素的前一个元素**,即为重复元素之后的第1个非重复元素.

- step4. 覆盖
	- 由于题目中给定的要求是**不需要考虑数组中超出新长度后面的元素**,所以直接用非重复元素覆盖重复元素的值即可

循环这个过程,直到遍历结束或**在重复元素之后找不到大于重复元素的前一个元素**时,结束循环.

此处要解释的是为什么要与重复元素的前一个元素比较,而非与重复元素比较.因为重复元素的前一个元素值是一个标识,该值标识了此时数组中已经被"去重"这个过程处理过的元素的最大值.后续"去重"操作的处理对象应该是大于重复元素的前一个元素值的元素.而非是大于重复元素的值的元素.

若在step2中找到的重复元素之后找不到大于重复元素的前一个元素的值,就说明给定的升序数组中,所有不重复的元素都已经被"去重"过程处理过了.

```

以 [0, 0, 1, 1, 1, 2, 2, 3, 3, 4] 为例.

nums[1]应该被nums[2]覆盖.

nums[3]应该被nums[5]覆盖.

nums[4]应该被nums[7]覆盖.若用非重复元素和重复元素比较大小,则会导致的后果是nums[4]被nums[6]覆盖

```

##### e. 实现

```go
package lc26_removeDuplicates

func removeDuplicates(nums []int) int {
	// 空数组或只有1个元素的数组 必然没有重复元素
	if len(nums) == 0 || len(nums) == 1 {
		return len(nums)
	}

	// 数组中至少有2个元素 此时第1个元素的值必然为非重复值
	var elementsCounter int = 1

	for i := 1; i < len(nums); i++ {
		// 后一个元素的值 > 前一个元素的值 则后一个元素必然是非重复元素
		if nums[i] > nums[i - 1] {
			elementsCounter++
		} else {
			// 后一个元素的值 <= 前一个元素的值 则后一个元素是重复元素
			// 从该元素向后 找到第1个 > 前一个元素值的元素 覆盖
			found := false
			for j := i + 1; j < len(nums); j++ {
				if nums[j] > nums[i - 1] {
					nums[i] = nums[j]
					elementsCounter++
					found = true
					break
				}
			}

			// 若向后没有找到 > 前一个元素值的元素 则说明后面所有元素都比前一个元素小
			// 说明后面都是重复元素了
			if !found {
				return elementsCounter
			}
		}
	}

	return elementsCounter
}

```

##### f. 算法改进-O(n)的算法

这个算法的时间复杂度是O(n^2).对每一个元素的遍历是O(n)的时间复杂度,从每一个元素向后找到大于该元素的前一个元素又是一个O(n)的时间复杂度.故为O(n^2).

而且这个算法强依赖于一个已知条件:数组中的元素**按升序排序**.也就是代码中第21行的判断条件,若题目条件改为数组中的元素按降序排序,则需要修改代码(将`>`修改为`<`).说明这个算法实现的抽象层级并不高.

但实际上这道题只需使用O(n)的时间复杂度即可解决.我们考虑"去重"这个操作,其核心含义为:**任取数组中的一个元素,若该元素已经在结果数组中出现过,则该元素不需要再出现在结果数组中;否则该元素需要出现在结果数组中.**

那么该如何判定某个元素是否已经在结果数组中出现过呢?此时就需要依靠题目中的条件:数组中的元素已经按升序**排序**.重点在于元素是有序的,而非是元素排列是按升序或降序的.元素有序,意味着值相同的元素在数组中是**连续出现**的,这也就成为了我们判定某个元素在数组中是否已经出现过的依据.

依旧要使用到双指针.其中一个指针(指针A)指向当前元素,用于判定当前元素是否为第1次出现的元素;另一个指针(指针B)指向用于覆盖的位置.而最终题目要求的结果数组长度,即为指针B位移的次数.当指针A移动到数组末尾时,去重操作就完成了.该算法的时间复杂度为O(n).

```go
package lc26_removeDuplicates

func removeDuplicates(nums []int) int {
	elementsCounter := 0

	for i := 0; i < len(nums); i++ {
		if i == 0 {
			elementsCounter++
		} else {
			if nums[i] != nums[i - 1] {
				// 若 nums[i] != nums[i - 1] 则说明nums[i]是第1次出现的
				// 需要放到结果数组中
				nums[elementsCounter] = nums[i]
				elementsCounter++
			}
		}
	}

	return elementsCounter
```

##### g. 代码上改进-条件合并

可以看到循环中有2个条件分支:`i == 0`和`nums[i] != nums[i - 1]`.但实际上这两个分支下,做的操作是完全相同的.虽然从代码上来看,`i == 0`时只有`elementsCounter++`的操作,但实际上是因为此时`elementsCounter`和`i`值相同,所以不需要`nums[elementsCounter] = nums[i]`这一步操作,换言之,`i == 0`是`nums[i] != nums[i - 1]`的特例.因为当`i == 0`时,数组中的第1个元素必然是不重复的,即必然要保留在结果数组中的.

即:**循环中的条件分支是可以合并的.**

```go
package lc26_removeDuplicates

func removeDuplicates(nums []int) int {
	elementsCounter := 0

	for i := 0; i < len(nums); i++ {
		if i == 0 || nums[i] != nums[i - 1] {
			nums[elementsCounter] = nums[i]
			elementsCounter++
		}
	}

	return elementsCounter
}
```

#### 1.3.3 移动零

##### a. 题目要求

给定一个数组`nums`,编写一个函数将所有`0`移动到数组的末尾,同时保持非零元素的相对顺序.

示例:

```

输入: [0,1,0,3,12]
输出: [1,3,12,0,0]

```

说明:

1. 必须在原数组上操作,不能拷贝额外的数组
2. 尽量减少操作次数

##### b. 审题

实际上本题的核心思路和上一题是一样的.本质上都是对数组中的任意元素,当该元素满足某些条件时,对该元素做某些操作.

上一题中的"某些条件"为:某元素为数组中重复出现的元素

上一题中的"某些操作"为:删除该元素

本题中的"某些条件"为:元素值为0

本题中的"某些操作"为:将该元素移动到数组末尾

那么解题的大方向即为:**遍历数组,将0放到末尾,非零元素不动**.也就是说我们最终设计出来的算法,其时间复杂度为O(n).

整体思路有了,接下来要解决的是细节问题

- 细节问题1:如何把0从数组中部移动到尾部?

实际上移动操作共2个子操作:

1. 将当前位置上的0删除
	- 若当前元素值为0,则从当前元素的后一个元素开始,直到原数组的末尾为止,向前位移1个单位长度.即**覆盖当前的0**
2. 在数组末尾用0覆盖之前的值

此处要注意的是,一定是**原数组的末尾**,因为后续还有用0覆盖的操作,所以向前位移的操作目标为原数组中从当前位置向后的所有元素,这其中不包含我们后续覆盖用的0.

- 细节问题2:在数组末尾用0覆盖时,具体应该放在哪里?

此时就需要一个指向当前的0应该覆盖的目标位置,每有一个0覆盖了其他元素,该位置向前位移1个单位长度.

##### c. 根据审题结果,寻找合适的数据结构

很明显本题还是使用双指针.一个指针(指针A)指向当前正在遍历的元素;另一个指针(指针B)指向用0覆盖的元素.

##### d. 解题思路

遍历数组,若当前元素值为0,则:

- step1. 将从当前位置向后的所有原数组中的元素向前位移1个单位
- step2. 将数组末尾处的值改为0
- step3. 指针B向前位移1个单位长度

注意:遍历数组时长度不能使用`len(nums)`,因为指针A指向0时,这个0会被后续元素覆盖掉,但覆盖后指针A指向的值是否为0此时是未知的,因此指针A的位置是不能改变的.如果使用`len(nums)`,则这个值始终不变,而我们需要的是让这个长度`-1`,因为已经有一个值为0的元素被处理过了,因此需要用一个变量记录初始状态数组的长度,用这个变量的值作为循环的次数,当遇到一个0时,该变量值-1.

##### e. 实现

```go
package lc283_moveZeroes

func moveZeroes(nums []int)  {
	originLen := len(nums)
	zeroTarget := len(nums) - 1

	for i := 0; i < originLen; i++ {
		if nums[i] == 0 {
			copy(nums[i: originLen], nums[i + 1: originLen])
			nums[zeroTarget] = 0
			i--
			zeroTarget--
			originLen--
		}
	}
}
```

##### f. 算法改进-O(n)的算法

注意这个算法的第9行.`copy()`操作的时间复杂度为O(len(nums) - i),也就是说这个操作的时间复杂度也是一个O(n).因此实际上这个算法的时间复杂度并不是O(n)而是O(n^2).

再来对比本题与上一题:

上一题中,当元素为重复值,则不要;本题中,当元素值为0,则不要.因此本题的核心思路为:**遍历数组,若元素值不为0,则将该元素放置到正确的位置上去.最后再用0覆盖数组末尾的值.**

那么就还需要一个指针指向"正确的位置".所谓正确的位置,即为从0开始,每放置了一个非0元素,则该指针向后位移1个单位.

```go
package lc283_moveZeroes

func moveZeros(nums []int) {
	nonzeroPosition := 0

	for i := 0; i < len(nums); i++ {
		// 元素值非0就要
		if nums[i] != 0 {
			nums[nonzeroPosition] = nums[i]
			nonzeroPosition++
		}
	}

	// 从nonzeroPosition处开始 直到数组结尾处 皆为要用0覆盖的位置
	for nonzeroPosition < len(nums) {
		nums[nonzeroPosition] = 0
		nonzeroPosition++
	}
}
```

该算法的时间复杂度为O(2n),即O(n).

## 2. 链表

链表是一个**有序集合**,该集合需要维护元素之间的相对位置,但并不需要在连续的内存空间中维护这些元素.

从这个定义中,可以获知2个关键点:

1. 链表在内存空间的分布是**非连续的**
2. 不能随机访问链表中的节点

### 2.1 单链表

#### 2.1.1 单链表的定义

链表中的节点只有指向其后一个节点的指针,没有指向其前一个节点的指针.也就是说,遍历链表时,只能访问到某个节点之后的节点,无法访问到某个节点之前的节点.

#### 2.1.2 单链表的实现

##### a. 节点的实现

```go
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
```

##### b. 单链表的实现

```go
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
	return l.head == nil
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
	tailNode := &Node{}
	tailNode.SetValue(value)
	// 细节:若LinkedList为空 直接设置Head指针指向待Append的Node即可
	if l.len == 0 {
		l.head = tailNode
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
		lastNode.SetNext(tailNode)
	}

	l.len++
}

func(l *LinkedList) Insert(index int, value interface{}) error {
	// 细节:判定index是否合法
	if index < 0 || index > l.len {
		return errors.New("illegal index")
	}

	// 细节:index == 0 等价于Prepend操作
	if index == 0 {
		l.Prepend(value)
		return nil
	}

	// 细节:index == len 等价于Append操作
	if index == l.len {
		l.Append(value)
		return nil
	}

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

	return nil
}

func (l *LinkedList) Delete(index int) error {
	// 细节:判定l是否是一个空Linked List
	if l.head == nil {
		return errors.New("can't delete element in a nil linked list")
	}

	// 细节:判定index的是否合法
	if index < 0 || index > l.len - 1 {
		return errors.New("illegal index")
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

	return nil
}
```

### 2.2 双链表

#### 2.2.1 双链表的定义

链表中的节点不仅有指向其后一个节点的指针,还有指向其前一个节点的指针.也就是说,遍历链表时,不仅能够访问到某个节点之后的节点,还能够访问到该节点之前的节点.

#### 2.2.2 双链表的实现

##### a. 节点的实现

```go
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
```

##### b. 双链表的实现

```go
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
```


	

















































