package run

import (
	"context"
	"io"
	"sync"
	"time"
)

var HealthyTimeoutKey = HealthyContextTimeout{}

type HealthyContextTimeout struct {
	T    time.Duration
	Rest <-chan bool
}

type Healthy struct {
	doneCtx context.Context
	cancel  context.CancelFunc
	once    *sync.Once
}

func NewHealthy() *Healthy {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Healthy{
		doneCtx: ctx,
		cancel:  cancel,
		once:    new(sync.Once),
	}
}

func (h *Healthy) loopWait(ctx context.Context, timeout time.Duration, rest <-chan bool) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			return
		case <-rest:
			ticker.Reset(timeout)
		}
	}
}

//Monitoring do not call repeatedly (only the first call is valid)
//if need to detect timeout, please set HealthyContextTimeout value for monitoring
func (h *Healthy) Monitoring(ctx context.Context, done func()) {
	h.once.Do(func() {
		defer h.cancel() //final execution. make sure all cleanup is done
		defer done()
		v := ctx.Value(HealthyTimeoutKey)
		if v != nil {
			switch ht := v.(type) {
			case HealthyContextTimeout:
				h.loopWait(ctx, ht.T, ht.Rest)
				return
			case *HealthyContextTimeout:
				h.loopWait(ctx, ht.T, ht.Rest)
				return
			default:
			}
		}
		<-ctx.Done()
	})
}

//DoneContext must be called after Monitoring called
func (h *Healthy) DoneContext() context.Context {
	return h.doneCtx
}

//NewHealthyCloser close Healthy
func NewHealthyCloser(healthy *Healthy, cancel context.CancelFunc) io.Closer {
	return NewOnceCloser(func() error {
		cancel()
		<-healthy.DoneContext().Done()
		return nil
	})
}
