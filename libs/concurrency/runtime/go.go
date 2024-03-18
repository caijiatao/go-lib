package runtime

func Go(fn func()) {
	go func() {
		defer HandleCrash()
		fn()
	}()
}
