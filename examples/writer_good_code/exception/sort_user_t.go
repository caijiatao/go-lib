package main

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

type Comparator[T any] func(a, b T) bool

func memoryUsage() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

func SortObjects[T any](items []T, compare Comparator[T]) {
	startMemory := memoryUsage() // 记录开始时的内存使用
	start := time.Now()

	sort.Slice(items, func(i, j int) bool {
		return compare(items[i], items[j])
	})

	elapsed := time.Since(start)
	endMemory := memoryUsage() // 记录结束时的内存使用

	// 打印排序耗时和内存使用变化
	fmt.Printf("Sorting took %s, Memory increased by %d bytes\n", elapsed, endMemory-startMemory)
}

func main() {
	users := []User{
		//...
	}

	SortObjects(users, sortByName)

	SortObjects(users, sortByTransactionNum)

	SortObjects(users, sortByLastTransactionDate)
}
