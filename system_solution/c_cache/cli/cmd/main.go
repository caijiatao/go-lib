package main

import (
	"bufio"
	"c_cache/cli"
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	client := cli.NewClient()

	// 检查连接是否成功
	ctx := context.Background()
	result, err := client.Ping(ctx)
	if err != nil {
		fmt.Println("连接 c_cache 失败:", err)
		return
	}
	fmt.Println("成功连接到 c_cache!")
	fmt.Println("输入 c_cache 命令，输入 'exit' 退出。")

	// 开始交互式命令行
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			fmt.Println("退出 c_cache CLI")
			break
		}

		// 解析命令和参数
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		// 执行命令
		cmd := args[0]
		cmdArgs := args[1:]
		result, err = client.Do(ctx, cmd, cmdArgs...)
		if err != nil {
			fmt.Println("执行命令出错:", err)
			continue
		}

		fmt.Printf("(结果) %v\n", result)
	}
}
