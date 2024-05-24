package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)

	// 获取硬盘大小，转换成TB显示
	d, _ := mem.SwapMemory()
	fmt.Printf("SwapMemory: %v, %v, %v\n", d.Total/1024/1024, d.Free, d.UsedPercent)

}
