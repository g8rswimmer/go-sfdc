package bulk

import (
	"context"
	"math/rand"
	"time"
)

// Backoff implements exponential backoff.
// The wait time between retries is a random value between 0 and the "retry envelope".
// The envelope starts at Initial and increases by the factor of Multiplier every retry,
// but is capped at Max.
type Backoff struct {
	// Initial is the initial value of the retry envelope, defaults to 1 second.
	Initial time.Duration

	// Max is the maximum value of the retry envelope, defaults to 30 seconds.
	Max time.Duration

	// Multiplier is the factor by which the retry envelope increases.
	// It should be greater than 1 and defaults to 2.
	Multiplier float64

	// cur is the current retry envelope
	cur time.Duration
}

// Pause returns the next time.Duration that the caller should use to backoff.
func (bo *Backoff) Pause() time.Duration {
	if bo.Initial == 0 {
		bo.Initial = time.Second
	}
	if bo.cur == 0 {
		bo.cur = bo.Initial
	}
	if bo.Max == 0 {
		bo.Max = 30 * time.Second
	}
	if bo.Multiplier < 1 {
		bo.Multiplier = 2
	}
	// Select a duration between 1ns and the current max. It might seem
	// counterintuitive to have so much jitter, but
	// https://www.awsarchitectureblog.com/2015/03/backoff.html argues that
	// that is the best strategy.
	d := time.Duration(1 + rand.Int63n(int64(bo.cur)))
	bo.cur = time.Duration(float64(bo.cur) * bo.Multiplier)
	if bo.cur > bo.Max {
		bo.cur = bo.Max
	}
	return d
}

// Retry ...
func Retry(bo Backoff, f func() (stop bool, err error)) error {
	for {
		stop, err := f()
		if stop || err != nil {
			return err
		}
		p := bo.Pause()
		time.Sleep(p)
	}
}

// SleepWithContext is similar to time.Sleep, but it can be interrupted by ctx.Done() closing.
func SleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// RetryWithContext ...
func RetryWithContext(ctx context.Context, bo Backoff, f func() (stop bool, err error),
	sleep func(context.Context, time.Duration) error) error {
	for {
		stop, err := f()
		if stop || err != nil {
			return err
		}
		p := bo.Pause()
		if cerr := sleep(ctx, p); cerr != nil {
			return cerr
		}
	}
}
