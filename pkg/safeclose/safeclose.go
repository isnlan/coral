package safeclose

type Close struct {
	C chan struct{}
}

func New() *Close {
	return &Close{C: make(chan struct{})}
}

func (c *Close) Close(f func() error) error {
	select {
	case <-c.C:
		return nil
	default:
		close(c.C)
		if f != nil {
			return f()
		}
		return nil
	}
}

func (c *Close) IsClosed() bool {
	select {
	case <-c.C:
		return true
	default:
		return false
	}
}
