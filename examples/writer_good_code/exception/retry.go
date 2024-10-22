package exception

type Message struct{}

func asyncRetrySendUserMessage(message Message) {}

func sendUserMessage(message Message) error {
	return nil
}

func CompleteOrder(orderID string) error {
	// 完成订单的其他逻辑...
	message := Message{}
	err := sendUserMessage(message)
	if err != nil {
		asyncRetrySendUserMessage(message)
	}
	return nil
}
