package compose

type ChainFunc[T any] func(T) T

// Usage:
// Turns `hf1(hf2(hf3(h)))` -> `compose.NewChain(hf1, hf2, hf3).Compose(h)`

type Chain[T any] struct {
	chain []ChainFunc[T]
}

func NewChain[T any](cfs ...ChainFunc[T]) *Chain[T] {
	chain := make([]ChainFunc[T], len(cfs))
	copy(chain, cfs)

	return &Chain[T]{
		chain: chain,
	}
}

// Next adds cf to the end of the running chain
func (c *Chain[T]) Next(cf ChainFunc[T]) *Chain[T] {
	c.chain = append(c.chain, cf)
	return c
}

// Compose resolves the chain into a single http.Handler
//
// If `end` is nil, defaults to http.DefaultServeMux
func (c *Chain[T]) Compose(arg T) T {
	r := arg
	for i := len(c.chain) - 1; i >= 0; i-- {
		r = c.chain[i](r)
	}

	return r
}
