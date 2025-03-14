package compose

import "slices"

// Usage:
// Turns `hf1(hf2(hf3(h)))` -> `compose.New(hf1, hf2, hf3).Compose(h)`

type ChainFunc[T any] = func(T) T

type Chain[T any] struct {
	chain []ChainFunc[T]
}

// New creates a new chain of passed in ChainFuncs.
//
// No ChainFunc functions are called until Compose is called.
func New[T any](cfs ...ChainFunc[T]) *Chain[T] {
	return &Chain[T]{
		chain: slices.Clone(cfs),
	}
}

// Next creates a new Chain with `cf` added to the previous chain
// leaving the original Chain unchanged
func (c *Chain[T]) Next(cf ChainFunc[T]) *Chain[T] {
	cc := New(c.chain...)
	cc.chain = append(cc.chain, cf)
	return cc
}

// Merge combines two chains into a single chain
// leaving the original chains untouched
func (c *Chain[T]) Merge(c2 *Chain[T]) *Chain[T] {
	cc := New(c.chain...)
	cc.chain = append(cc.chain, c2.chain...)
	return cc
}

// Compose resolves the chain into a single T value
func (c *Chain[T]) Compose(arg T) T {
	r := arg
	for i := len(c.chain) - 1; i >= 0; i-- {
		r = c.chain[i](r)
	}

	return r
}
