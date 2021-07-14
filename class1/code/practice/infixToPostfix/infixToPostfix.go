package infixToPostfix

import (
	stack2 "code/dataStructure/stack"
)

var operatorPriority map[string]int = map[string]int{
	"+": 0,
	"-": 0,
	"*": 1,
	"/": 1,
}

func InfixToPostfix(infix string) (postfix string) {
	var operatorStack stack2.Stack = stack2.Stack{}
	for i := 0; i < len(infix); i++ {
		_, ok := operatorPriority[string(infix[i])]
		if !ok {
			// 操作数 直接放入结果
			postfix += string(infix[i])
			continue
		}

		// 操作符
		for {
			// 栈为空 将当前操作符压入栈内
			if operatorStack.IsEmpty() {
				operatorStack.Push(infix[i])
				break
			}

			// 栈不为空 检查栈顶操作符的优先级 与当前操作符的优先级比对
			topOperator, _ := operatorStack.Peek()

			// 栈顶操作符优先级 < 当前操作符优先级 将当前操作符压入栈内
			if operatorPriority[string(topOperator.(uint8))] < operatorPriority[string(infix[i])] {
				operatorStack.Push(infix[i])
				break
			}

			// 栈顶操作符优先级 >= 当前操作符优先级 则弹出栈顶操作符 放入结果中
			top, _ := operatorStack.Pop()
			postfix += string(top.(uint8))
		}
	}

	for !operatorStack.IsEmpty() {
		topOperator, _ := operatorStack.Pop()
		postfix += string(topOperator.(uint8))
	}
	return
}

func InfixToPostfixWithParentheses(infix string) (postfix string) {
	var operatorStack stack2.Stack = stack2.Stack{}
	for i := 0; i < len(infix); i++ {
		// 左括号 入栈
		if string(infix[i]) == "(" {
			operatorStack.Push(infix[i])
			continue
		}

		_, ok := operatorPriority[string(infix[i])]

		// 操作数 直接放入结果
		if string(infix[i]) != "(" && string(infix[i]) != ")" && !ok {
			postfix += string(infix[i])
			continue
		}

		// 右括号 持续出栈直到出栈元素为左括号为止
		if string(infix[i]) == ")" {
			for {
				topOperator, _ := operatorStack.Pop()
				if string(topOperator.(uint8)) == "(" {
					break
				}

				postfix += string(topOperator.(uint8))
			}
			continue
		}

		// 操作符
		for {
			// 栈为空 将当前操作符压入栈内
			if operatorStack.IsEmpty() {
				operatorStack.Push(infix[i])
				break
			}

			// 栈不为空 检查栈顶操作符的优先级 与当前操作符的优先级比对
			topOperator, _ := operatorStack.Peek()

			// 栈顶操作符优先级 < 当前操作符优先级 或 栈顶为(时 将当前操作符压入栈内
			topStr := string(topOperator.(uint8))
			if operatorPriority[topStr] < operatorPriority[string(infix[i])] || topStr == "(" {
				operatorStack.Push(infix[i])
				break
			}

			// 栈顶操作符优先级 >= 当前操作符优先级 则弹出栈顶操作符 放入结果中
			top, _ := operatorStack.Pop()
			postfix += string(top.(uint8))
		}
	}

	for !operatorStack.IsEmpty() {
		topOperator, _ := operatorStack.Pop()
		postfix += string(topOperator.(uint8))
	}
	return
}