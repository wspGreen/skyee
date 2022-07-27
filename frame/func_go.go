package frame

func Go(fn func()) {
	go func() {
		fn()
	}()
}
