package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := ""                // 替换为你要测试的API的URL
	numRequests := 100       // 替换为你要发送的请求数量
	concurrentRequests := 10 // 替换为你想要的并发请求数量

	var wg sync.WaitGroup
	requests := make(chan struct{}, concurrentRequests)

	startTime := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-requests:
					sendRequest(url)
				}
			}
		}()
	}

	for i := 0; i < numRequests; i++ {
		requests <- struct{}{}
	}

	close(requests)

	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Completed %d requests in %s\n", numRequests, elapsedTime)
}

func sendRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 在这里处理响应，如果需要的话
	// 例如，你可以读取响应体并进行验证

	// 例如：
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//     fmt.Println("Error reading response body:", err)
	//     return
	// }
	// fmt.Println("Response:", string(body))
}
