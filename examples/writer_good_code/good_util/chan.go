package good_util

import (
	"github.com/samber/lo"
	"sync"
)

func loChannelDispatcher() {
	ch := make(chan int, 20)
	for i := 0; i <= 15; i++ {
		ch <- i
	}

	dispatchCount := 5     // 确定需要拆分的channel数量
	channelBufferCap := 10 // 拆分出来channel的缓冲区大小
	children := lo.ChannelDispatcher(ch, dispatchCount, channelBufferCap, lo.DispatchingStrategyRoundRobin[int])
	// children : []<-chan int{...}

	var wg sync.WaitGroup
	wg.Add(dispatchCount)

	// 编写消费者函数，这里根据业务逻辑进行改写。
	consumer := func(c <-chan int) {
		for {
			msg, ok := <-c
			if !ok {
				println("closed")
				wg.Done()
				break
			}

			println(msg)
		}
	}

	// 启动多个协程对拆分出来的children channel进行消费
	for i := range children {
		go consumer(children[i])
	}

	close(ch)

	wg.Wait()
}

func loSliceToChannel() {
	list := []int{1, 2, 3, 4, 5}

	for v := range lo.SliceToChannel(2, list) {
		println(v)
	}
}

func loBuffer() {
	ch := lo.SliceToChannel(2, []int{1, 2, 3, 4, 5})

	items1, length1, duration1, ok1 := lo.Buffer(ch, 3)
	// []int{1, 2, 3}, 3, 0s, true
	println(items1, length1, duration1, ok1)

	items2, length2, duration2, ok2 := lo.Buffer(ch, 3)
	// []int{4, 5}, 2, 0s, false
	println(items2, length2, duration2, ok2)
}
