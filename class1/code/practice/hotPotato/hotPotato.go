package hotPotato

import "code/dataStructure/queue"

var players []string = []string{"Bill", "David", "Susan", "Jane", "Kent", "Brad"}

func HotPotato(nums int) (winner string) {
	var playersQueue queue.Queue = queue.Queue{}
	for _, player := range players {
		playersQueue.Enqueue(player)
	}

	for playersQueue.Size() > 1 {
		for i := 0; i < nums; i++ {
			handOver, _ := playersQueue.Dequeue()
			playersQueue.Enqueue(handOver)
		}
		// 传递次数到达指定数量 移除持有土豆的玩家 并将土豆的持有权顺延给下一位玩家
		playersQueue.Dequeue()
	}

	winPlayer, _ := playersQueue.Dequeue()
	return winPlayer.(string)
}
