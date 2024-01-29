package rule_dig

import (
	"fmt"
)

func demoDig() {
	// 定义矩阵
	matrix := [][]int{
		{100, 50, 23, 27, 10, 97, 3, 50, 40},
		{99, 49, 24, 26, 10, 99, 1, 40, 50},
		{105, 90, 10, 5, 10, 100, 0, 45, 45},
	}

	// 检查规则 A1 = A2 + A3
	for row := 0; row < len(matrix); row++ {

		rule1Satisfied := checkRule1(matrix, row)

		// 检查规则 A5 = A6 + A7 - A8 - A9
		rule2Satisfied := checkRule2(matrix, row)

		// 输出结果
		fmt.Println("规则 A1 = A2 + A3 是否满足:", rule1Satisfied)
		fmt.Println("规则 A5 = A6 + A7 - A8 - A9 是否满足:", rule2Satisfied)
	}
}

// 检查规则 A1 = A2 + A3
func checkRule1(matrix [][]int, row int) bool {
	return matrix[row][0] == matrix[row][1]+matrix[row][2]
}

// 检查规则 A5 = A6 + A7 - A8 - A9
func checkRule2(matrix [][]int, row int) bool {
	return matrix[row][4] == matrix[row][5]+matrix[row][6]-matrix[row][7]-matrix[row][8]
}
