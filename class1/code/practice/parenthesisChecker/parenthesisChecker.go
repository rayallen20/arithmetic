package parenthesisChecker

import (
	stack2 "code/dataStructure/stack"
	"strings"
)

var leftSymbols string = "([{"

var rightSymbols string = ")]}"

func ParenthesisChecker(symbols string) bool {
	symbolRunes := []rune(symbols)
	stack := stack2.Stack{}
	for i := 0; i < len(symbolRunes); i++ {
		// 左括号 入栈
		if strings.ContainsRune(leftSymbols, symbolRunes[i]) {
			stack.Push(string(symbolRunes[i]))
		}

		// 右括号 栈顶比对
		if strings.ContainsRune(rightSymbols, symbolRunes[i]) {
			top, _ := stack.Peek()
			if matches(top.(string), string(symbolRunes[i])) {
				stack.Pop()
			} else {
				return false
			}
		}
	}

	return stack.IsEmpty()
}

func matches(left, right string) bool {
	return strings.Index(leftSymbols, left) == strings.Index(rightSymbols, right)
}
