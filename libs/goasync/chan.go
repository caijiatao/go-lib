package goasync

// SafeClose 安全地关闭一个 channel，避免因重复关闭或关闭 nil channel 而导致的 panic。
// 如果 channel 为 nil，返回 false；如果成功关闭，返回 true。
// 使用示例：
//   ch := make(chan struct{})
//   if SafeClose(ch) {
//       // channel 已成功关闭
//   }
func SafeClose(ch chan struct{}) (ok bool) {
	if ch == nil {
		return false
	}

	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	close(ch)
	return true
}
