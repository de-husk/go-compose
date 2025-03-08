package compose

import "slices"

// Usage:
// Turns `hf1(hf2(hf3(h)))` -> `compose.NewChain(hf1, hf2, hf3).Compose(h)`

type ChainFunc[T any] func(T) T

type Chain[T any] struct {
	chain []ChainFunc[T]
}

// NewChain creates a new chain of passed in ChainFuncs.
//
// No ChainFunc functions are called until Compose is called.
func NewChain[T any](cfs ...ChainFunc[T]) *Chain[T] {
	return &Chain[T]{
		chain: slices.Clone(cfs),
	}
}

// Next adds cf to the end of the running chain
func (c *Chain[T]) Next(cf ChainFunc[T]) *Chain[T] {
	c.chain = append(c.chain, cf)
	return c
}

// Compose resolves the chain into a single T value
func (c *Chain[T]) Compose(arg T) T {
	r := arg
	for i := len(c.chain) - 1; i >= 0; i-- {
		r = c.chain[i](r)
	}

	return r
}
