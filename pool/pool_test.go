package pool_test

import (
	"context"
	"testing"
	"time"

	"github.com/ic-n/preallocated-pool/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Value struct {
	Index    int
	Duration time.Duration
}

func TestPool(t *testing.T) {
	ctx := context.Background()
	start := time.Now()
	c := make(chan Value, 4)

	p := pool.New[Value](2, time.Second)
	for i := 0; i < 4; i++ {
		err := p.With(ctx, func(v *Value) {
			v.Index = i
			v.Duration = time.Since(start)
			c <- *v
		})
		require.NoError(t, err)
	}

	require.GreaterOrEqual(t, time.Second, (<-c).Duration)
	require.GreaterOrEqual(t, time.Second, (<-c).Duration)
	require.LessOrEqual(t, time.Second, (<-c).Duration)
	require.LessOrEqual(t, time.Second, (<-c).Duration)
}

func TestPoolWithCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()

	p := pool.New[Value](2, time.Second)
	for i := 0; i < 2; i++ {
		err := p.With(ctx, func(v *Value) {
			v.Index = i
			v.Duration = time.Since(start)
		})
		require.NoError(t, err)
	}
	for i := 0; i < 2; i++ {
		err := p.With(ctx, func(v *Value) {
			v.Index = i
			v.Duration = time.Since(start)
		})
		if assert.Error(t, err) {
			assert.EqualError(t, err, "rate: Wait(n=1) would exceed context deadline")
		}
	}
}
