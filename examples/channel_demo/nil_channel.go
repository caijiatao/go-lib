package channel_demo

import (
	"fmt"
	"time"
)

func NilChannel() {
	c1, c2 := make(chan int), make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		c1 <- 1
		close(c1)
	}()
	go func() {
		time.Sleep(3 * time.Second)
		c2 <- 2
		close(c2)
	}()

	var ok1, ok2 bool
	for {
		select {
		case x := <-c1:
			ok1 = true
			fmt.Println("c1:", x)
			c1 = nil
		case x := <-c2:
			ok2 = true
			fmt.Println("c2:", x)
			c2 = nil
		}
		if ok1 && ok2 {
			break
		}
	}

	fmt.Println("done")
}
