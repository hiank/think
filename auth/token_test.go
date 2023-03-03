package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/hiank/think/auth"
	"gotest.tools/v3/assert"
)

func TestToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tk := auth.Export_newToken(ctx, auth.WithTokenTimeout(time.Microsecond))
	assert.Equal(t, tk.Err(), nil)
	assert.Equal(t, tk.ToString(), "")
	tk2 := tk.Fork(auth.WithTokenTimeout(time.Second))
	assert.Equal(t, tk2.Err(), nil)
	<-time.After(time.Millisecond * 20)
	assert.Equal(t, tk.Err(), context.DeadlineExceeded)
	assert.Equal(t, tk2.Err(), context.DeadlineExceeded)
	///

	tk = auth.Export_newToken(context.WithValue(ctx, auth.Export_contextkeyToken, "tt-key"))
	tk2 = tk.Fork(auth.WithTokenTimeout(time.Microsecond))
	assert.Equal(t, tk2.ToString(), "tt-key")
	<-time.After(time.Millisecond * 20)
	assert.Equal(t, tk.Err(), nil)
	assert.Equal(t, tk2.Err(), context.DeadlineExceeded)

	tk.Close()
	assert.Equal(t, tk.Err(), context.Canceled)
}

func TestTokenset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := auth.NewTokenset(ctx)
	tk := ts.Derive("110", auth.WithTokenTimeout(time.Microsecond))
	assert.Equal(t, tk.ToString(), "110")
	assert.Equal(t, tk.Err(), nil)

	tk2 := ts.Derive("110", auth.WithTokenTimeout(time.Millisecond*100))
	tk3 := tk.Fork(auth.WithTokenTimeout(time.Millisecond * 100))

	<-time.After(time.Millisecond * 20)
	assert.Equal(t, tk.Err(), context.DeadlineExceeded)
	assert.Equal(t, tk3.Err(), tk.Err())
	assert.Equal(t, tk2.Err(), nil)

	err := ts.Kill("110")
	assert.Equal(t, err, nil)
	assert.Equal(t, tk2.Err(), context.Canceled)

	err = ts.Kill("110")
	assert.Equal(t, err, auth.ErrNonRootoken)

	// m := auth.Export_tokenSetM(ts)
	// assert.Equal(t, m.)
	ts.Derive("120")
	m := auth.Export_tokenSetM(ts)
	_, found := m.Load("120")
	assert.Assert(t, found)

	ts.Close()
	m = auth.Export_tokenSetM(ts)
	_, found = m.Load("120")
	assert.Assert(t, !found)

	// ctx, cancel = context.WithCancel(ctx)
	// cancel()
	// ts = auth.NewTokenset(ctx)
	// assert.Equal(t, ts.
}
