package pool

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Pool[T any] struct {
	Pool    sync.Pool
	limiter *rate.Limiter
}

func New[T any](size int, interval time.Duration) *Pool[T] {
	rl := rate.Every(interval)
	limiter := rate.NewLimiter(rl, size)
	p := Pool[T]{
		Pool:    sync.Pool{},
		limiter: limiter,
	}

	for i := 0; i < size; i++ {
		var v T
		p.Pool.Put(&v)
	}

	return &p
}

func (p *Pool[T]) With(ctx context.Context, fn func(*T)) error {
	if err := p.limiter.Wait(ctx); err != nil {
		return err
	}

	v := p.Pool.Get().(*T)
	fn(v)
	p.Pool.Put(v)

	return nil
}
