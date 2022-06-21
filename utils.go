package main

func copyChannel[T any](in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T, 1)
	out2 := make(chan T, 1)

	go func() {
		defer close(out1)
		defer close(out2)

		for event := range in {
			out1 <- event
			out2 <- event
		}
	}()

	return out1, out2
}
