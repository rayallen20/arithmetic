package palChecker

import (
	"code/dataStructure/deque"
)

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

func PalCheckerByDeque(str string) (isPal bool) {
	strRunes := []rune(str)
	var charDeque deque.Deque = deque.Deque{}
	for i := 0; i < len(strRunes); i++ {
		charDeque.AddRear(strRunes[i])
	}

	for charDeque.Size() > 1 {
		headChar, _ := charDeque.RemoveFront()
		tailChar, _ := charDeque.RemoveRear()
		if headChar != tailChar {
			return false
		}
	}
	return true
}