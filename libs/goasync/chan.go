package goasync

func SafeClose(ch chan struct{}) (ok bool) {
	defer func() {
		// ignore error
		_ = recover()
	}()
	close(ch)
	return true
}
