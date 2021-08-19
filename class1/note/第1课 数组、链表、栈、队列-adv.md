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




















































