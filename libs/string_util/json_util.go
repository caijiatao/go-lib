package string_util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadJsonFromFile(filePath string) []map[string]interface{} {
	// 打开 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	value, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	data := make([]map[string]interface{}, 0)
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil
	}
	return data
}

func CountNumber(nums []int, n int) (count int) {
	for i := 0; i < len(nums); i++ {
		// 如果要赋值 则 v := nums[i]
		if nums[i] == n {
			count++
		}
	}
	return
}

func CountNumberBad(nums []int, n int) (count int) {
	for index := 0; index < len(nums); index++ {
		value := nums[index]
		if value == n {
			count++
		}
	}
	return
}
