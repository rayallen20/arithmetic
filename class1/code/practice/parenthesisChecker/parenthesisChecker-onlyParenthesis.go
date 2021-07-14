package parenthesisChecker

import stack2 "code/dataStructure/stack"

func OnlyParenthesisChecker(symbol string) bool {
	stack := stack2.Stack{}
	for i := 0; i < len(symbol); i++ {
		if symbol[i] == '(' {
			stack.Push(symbol[i])
		}

		if symbol[i] == ')' {
			if stack.IsEmpty() {
				return false
			} else {
				stack.Pop()
			}
		}
	}

	return stack.IsEmpty()
}
