package string_util

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestReadJsonFromFile(t *testing.T) {
	data := ReadJsonFromFile("Explore-logs-2024-05-11 09_53_54.json")
	timeLen := len("1714266618")
	statistic := make(map[string]int)
	sortKey := make([]string, 0)
	for _, d := range data {
		d["timestamp"] = d["timestamp"].(string)[:timeLen]
		parseInt, err := strconv.ParseInt(d["timestamp"].(string), 10, 64)
		if err != nil {
			return
		}
		dTime := time.Unix(parseInt, 0)
		statistic[dTime.Format("2006-01-02")]++
		if statistic[dTime.Format("2006-01-02")] == 1 {
			sortKey = append(sortKey, dTime.Format("2006-01-02"))
		}
	}
	for _, k := range sortKey {
		fmt.Println(k, statistic[k])
	}
}

func generate(ch chan int) {
	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch)
}

func filter(ch chan int) {
	for {
		c, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(c)
	}
}

func TestPrime(t *testing.T) {
	ch := make(chan int, 10)
	go generate(ch)
	for i := 0; i < 10; i++ {
		go filter(ch)
	}

	time.Sleep(5 * time.Second)
}

func TestName(t *testing.T) {
	if a := 1; false {

	} else if b := 2; false {

	} else if c := 3; false {

	} else {
		fmt.Println(a, b, c)
	}
}
