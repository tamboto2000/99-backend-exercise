package gocb

type noopCb struct {
	name string
}

func NoopCB() Circuit {
	return noopCb{}
}

func (n noopCb) Name() string {
	return n.name
}

func (n noopCb) Status() CircuitStatus {
	return Closed
}

func (n noopCb) Run(f func() error) error {
	return f()
}
