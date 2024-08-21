package main

import (
	"fmt"
	"sort"
	"time"
)

// User 结构体，包含用户的相关信息和行为
type User struct {
	Name                string
	TransactionNum      int
	LastTransactionDate time.Time
}

// 定义手续费区间及其对应规则
var feeBrackets = []struct {
	MaxTransactions   int
	FeePerTransaction int
}{
	{5, 20},
	{20, 10},
	{1<<31 - 1, 1}, // 无上限的最大值，交易笔数 > 20时的手续费
}

// User 类型的方法，计算该用户的手续费
func (u *User) CalculateFee() int {
	totalFee := 0
	remainingTransactions := u.TransactionNum

	for _, bracket := range feeBrackets {
		transactions := min(remainingTransactions, bracket.MaxTransactions)
		totalFee += transactions * bracket.FeePerTransaction
		remainingTransactions -= transactions
		if remainingTransactions <= 0 {
			break
		}
	}
	return totalFee
}

// min 函数：获取较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Users 类型：一个 User 切片，具备排序行为
type Users []User

// 根据不同的排序方式排序
func (users Users) SortBy(sortBy string) {
	sortFuncMap := map[string]func(i, j int) bool{
		"name":                func(i, j int) bool { return users[i].Name < users[j].Name },
		"transactionNum":      func(i, j int) bool { return users[i].TransactionNum < users[j].TransactionNum },
		"lastTransactionDate": func(i, j int) bool { return users[i].LastTransactionDate.Before(users[j].LastTransactionDate) },
	}

	sort.Slice(users, sortFuncMap[sortBy])
}

func main() {
	// 初始化用户数据
	users := Users{
		{"张三", 4, mustParseTime("2006.1.2", "2024.2.24")},
		{"李四", 10, mustParseTime("2006.1.2", "2024.3.1")},
		{"王五", 50, mustParseTime("2006.1.2", "2024.9.21")},
	}

	// 计算每个人的手续费
	for _, user := range users {
		fee := user.CalculateFee()
		fmt.Printf("%s 的手续费是 %d 元\n", user.Name, fee)
	}

	// 按交易笔数排序
	users.SortBy("transactionNum")
	fmt.Println("\n按交易笔数排序后的名字：")
	for _, user := range users {
		fmt.Println(user.Name)
	}

	// 按交易日期排序
	users.SortBy("lastTransactionDate")
	fmt.Println("\n按交易日期排序后的名字：")
	for _, user := range users {
		fmt.Println(user.Name)
	}
}

// 辅助函数：将日期字符串解析为 time.Time
func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
