package divideBy2

import (
	stack2 "code/dataStructure/stack"
	"strconv"
)

func DivideBy2(dec int) (bin string) {
	var remainders stack2.Stack = stack2.Stack{}

	for dec > 1 {
		remainders.Push(dec % 2)
		dec = dec / 2
	}
	remainders.Push(dec)

	for !remainders.IsEmpty() {
		binChar, _ := remainders.Pop()
		bin += strconv.Itoa(binChar.(int))
	}
	return bin
}
