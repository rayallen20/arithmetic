package array

import (
	"errors"
	"fmt"
)

type Array struct {
	array []interface{}
}

func (a *Array) Insert(index int, value interface{})  {
	a.check()
	if  0 < index && index < len(a.array) {
		// 主体思路:向数组中(非头尾)插入
		a.array = append(a.array, nil)
		copy(a.array[index:],a.array[index - 1:])
		a.array[index] = value
	} else if index == 0 {
		// 细节:向头部插入
		a.Prepend(value)
	} else if index == len(a.array) {
		// 细节:向尾部插入
		a.Append(value)
	} else {
		// 细节:待插入索引比当前数组长度要大
		padding := make([]interface{}, index - len(a.array))
		a.array = append(append(a.array, padding...),value)
	}
}

func (a *Array) check() {
	if a.array == nil {
		a.array = make([]interface{}, 0, 100)
	}
}

func (a *Array) Prepend(value interface{}) {
	a.check()
	a.array = append(a.array, nil)
	copy(a.array[1:], a.array[:])
	a.array[0] = value
}

func (a *Array) Append(value interface{}) {
	a.check()
	a.array = append(a.array, value)
}

func (a *Array) Delete(index int) {
	if 0 < index && index < len(a.array) {
		// 主体思路:删除数组中的元素
		// 细节:删除尾部的元素(和主体思路一致,合并处理)
		a.array = append(a.array[:index], a.array[index + 1:]...)
	} else if index == 0 {
		// 细节:删除头部的元素
		a.array = a.array[1:]
	} else if 0 > index || index >= len(a.array) {
		// 细节:索引越界
		err := errors.New("out of range")
		fmt.Printf("%v\n", err.Error())
	}
}

func (a *Array) Lookup() {
	a.check()
	fmt.Printf("%v\n", a.array)
}
