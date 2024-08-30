package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan<- string) {
	time.Sleep(time.Duration(id) * time.Second)
	result := fmt.Sprintf("Worker %d done", id)
	ch <- result
}

func processor(ch <-chan string, done chan<- bool) {
	for result := range ch {
		// Process the result
		fmt.Printf("Processed result: %s\n", result)
	}
	done <- true
}

func main() {
	resultChan := make(chan string)
	doneChan := make(chan bool)

	for i := 1; i <= 5; i++ {
		go worker(i, resultChan)
	}

	go processor(resultChan, doneChan)

	time.Sleep(6 * time.Second)

	close(resultChan)

	<-doneChan
}
