package or

func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan interface{})
	go func() {
		defer close(out)
		select {
		case <-channels[0]:
		case <-Or(channels[1:]...):
		}
	}()
	return out
}
