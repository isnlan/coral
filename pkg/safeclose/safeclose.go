package safeclose

type Closer interface {
	Close(f func() error) error
	IsClosed() bool
	C() <-chan struct{}
}

type Close struct {
	c chan struct{}
}

func New() *Close {
	return &Close{c: make(chan struct{})}
}

func (c *Close) Close(f func() error) error {
	select {
	case <-c.c:
		return nil
	default:
		close(c.c)
		if f != nil {
			return f()
		}
		return nil
	}
}

func (c *Close) IsClosed() bool {
	select {
	case <-c.c:
		return true
	default:
		return false
	}
}

func (c *Close) C() <-chan struct{} {
	return c.c
}
