package postfixToInfix

import stack2 "code/dataStructure/stack"

func PostfixToInfix(postfix string) (infix string) {
	var operandStack stack2.Stack = stack2.Stack{}
	for i := 0; i < len(postfix); i++ {
		token := string(postfix[i])

		if token == " " {
			continue
		}

		operatorRes := isOperator(token)
		// 操作数 入栈
		if !operatorRes {
			operandStack.Push(token)
		} else {
			// 操作符 取出栈顶的2个元素 先弹出的元素放在操作符后面 后弹出的元素放在操作符前面
			// 将运算结果再压入栈内
			tailOperand, _ := operandStack.Pop()
			headOperand, _ := operandStack.Pop()
			operateRes := "(" + headOperand.(string) + token + tailOperand.(string) + ")"
			operandStack.Push(operateRes)
		}
	}
	res, _ := operandStack.Pop()
	return res.(string)
}

func isOperator(token string) bool {
	var operators []string = []string{"+", "-", "*", "/"}
	for _, v := range operators {
		if v == token {
			return true
		}
	}
	return false
}