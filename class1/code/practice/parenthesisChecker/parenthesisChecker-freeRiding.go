package parenthesisChecker

func FreeRidingParenthesisChecker(symbol string) bool {
	var leftParenthesis int
	var rightParenthesis int
	for i := 0; i < len(symbol); i++ {
		if symbol[i] == '(' {
			leftParenthesis++
		} else {
			rightParenthesis++
		}
	}

	return leftParenthesis == rightParenthesis
}
