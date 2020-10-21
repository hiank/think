package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/hiank/think/core"
)

func TestWithPort(t *testing.T) {

	assert.Equal(t, core.WithPort("192.168.1.22", 1024), "192.168.1.22:1024")
}

func TestHealthLock(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Millisecond)

	err := core.HealthLock(ctx, ticker, func() error {
		return errors.New("testClose")
	})
	assert.Equal(t, err.Error(), "testClose", "指定时间后将调用指定关闭方法")

	go func() {

		ticker2, ticker3 := time.NewTicker(time.Millisecond*500), time.NewTicker(time.Second*2)
		for {
			select {
			case <-ticker2.C:
				ticker.Reset(time.Second)
				ticker2.Reset(time.Millisecond * 500)
			case <-ticker3.C:
				cancel()
				return
			}
		}
	}()

	ticker.Reset(time.Second)
	err = core.HealthLock(ctx, ticker, func() error {
		return errors.New("testClose2")
	})

	select {
	case <-ctx.Done():
	default:
		assert.Assert(t, false, "必须是响应的context cancel")
	}

	select {
	case <-ticker.C:
		assert.Assert(t, false, "ticker必须未响应")
	default:
	}
	assert.Equal(t, err.Error(), "testClose2")
}
