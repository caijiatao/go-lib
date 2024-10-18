package good_t

// 定义手续费区间及其对应规则
var feeBrackets = []struct {
	MaxTransactions   int
	FeePerTransaction int
}{
	{5, 20},
	{15, 10},
	{1<<31 - 1, 1}, // 无上限的最大值，交易笔数 > 20时的手续费
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 根据交易笔数计算手续费
func calculateFee(transactionNum int) (totalFee int) {
	for _, bracket := range feeBrackets {
		// 获取区间交易笔数
		transactions := min(transactionNum, bracket.MaxTransactions)
		totalFee += transactions * bracket.FeePerTransaction

		// 已经计算的交易笔数可以减掉，等到全部计算完成则返回结果
		transactionNum -= transactions
		if transactionNum <= 0 {
			break
		}
	}
	return totalFee
}
