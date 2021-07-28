package palChecker

func  PalChecker(str string) (isPal bool) {
	strRunes := []rune(str)
	for i := 0; i < len(strRunes) / 2; i++ {
		head := strRunes[i]
		tail := strRunes[len(strRunes) - 1 - i]
		if head != tail {
			return false
		}
	}
	return true
}
