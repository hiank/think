package run

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/hiank/think/exp/lists"
)

const (
	ErrTimestampsetNotReady = Err("run: Timestampset not ready")
)

type Timestampset struct {
	// base time.Duration
	pp chan<- *Timestamp
}

func NewTimestampset(ctx context.Context, interval time.Duration) *Timestampset {
	pp := make(chan *Timestamp, 8)
	tss := &Timestampset{pp: pp}
	go tss.loop(ctx, pp, interval)
	return tss
}

//currentime 当前时间(目前使用的是系统当前时间，后续可以设计为某个时刻的相对时间)
func (tss *Timestampset) currentime() time.Duration {
	return time.Duration(time.Now().UnixMicro())
}

func (tss *Timestampset) loop(ctx context.Context, pp <-chan *Timestamp, interval time.Duration) {
	ticker, tsList, done := time.NewTicker(interval), list.New(), make(chan *list.Element, 8)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case elm := <-done:
			tsList.Remove(elm)
		case ts := <-pp:
			ts.listeners, ts.done, ts.elm = list.New(), done, lists.InsertBeforeFunc(tsList, ts, func(cur, want *Timestamp) bool {
				return cur.Cutime() >= want.Cutime()
			})
		case <-ticker.C:
			///
			currentime := tss.currentime()
			lists.Foreach(tsList, func(ts *Timestamp) (done bool) {
				ts.Step(currentime)
				return
			})
		}
	}
}

func (tss *Timestampset) New(duration time.Duration) (ts *Timestamp, err error) {
	//
	ts = &Timestamp{cutime: tss.currentime() + duration}
	select {
	case tss.pp <- ts:
	default:
		ts, err = nil, ErrTimestampsetNotReady
	}
	return
}

// func (tss *Timestampset)

type Timestamp struct {
	cutime    time.Duration
	elm       *list.Element
	done      chan<- *list.Element
	listeners *list.List
	mux       sync.RWMutex
}

func (ts *Timestamp) Step(currentime time.Duration) {
	// ts.mux.Lock()
	// defer ts.mux.Unlock()
	//
	cutime := ts.Cutime() ///the Cutime func would locked by 'ts.mux.Lock()'
	///
	ts.mux.Lock()
	defer ts.mux.Unlock()
	if cutime < currentime {
		if ts.elm != nil && ts.done != nil {
			ts.done <- ts.elm //notice tss for remove the elm
			ts.done, ts.elm = nil, nil
		}
	}
	lists.DeleteFunc(ts.listeners, func(listener *TimestampListener) (ok bool) {
		if ok = listener.Cutime <= currentime; ok {
			listener.Response(ts)
		}
		return
	})
}

//Listen 监听不超过指定剩余时间(此时刻第一次出现时，执行响应(仅一次))
//NOTE:
func (ts *Timestamp) Listen(left time.Duration, rsp func(*Timestamp)) {
	ts.mux.Lock()
	defer ts.mux.Unlock()
	listener := &TimestampListener{
		Cutime:   ts.Cutime() - left,
		Response: rsp,
	}
	lists.InsertBeforeFunc(ts.listeners, listener, func(cur, want *TimestampListener) bool {
		return cur.Cutime >= want.Cutime
	})
}

//Cutime cut time 截止时间(正常计时下，到此时刻触发结束响应)
func (ts *Timestamp) Cutime() time.Duration {
	ts.mux.RLock()
	defer ts.mux.RUnlock()
	return ts.cutime
}

//Speedup 加速，使截止时刻提前指定时长
func (ts *Timestamp) Speedup(dt time.Duration) {
	ts.mux.Lock()
	defer ts.mux.Unlock()
	ts.cutime -= dt
}

type TimestampListener struct {
	//
	Cutime   time.Duration //触发时刻
	Response func(*Timestamp)
}
