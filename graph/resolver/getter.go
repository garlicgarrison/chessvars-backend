package resolver

import (
	"context"
	"sync"
)

type getter[returnType any, fnType func(context.Context) (returnType, error)] struct {
	fn fnType

	once sync.Once

	reply returnType
	err   error
}

func NewGetter[R any, F func(context.Context) (R, error)](fn F) getter[R, F] {
	return getter[R, F]{
		fn: fn,
	}
}

func (g *getter[R, F]) Call(ctx context.Context) (R, error) {
	g.once.Do(func() {
		g.reply, g.err = g.fn(ctx)
	})
	return g.reply, g.err
}
