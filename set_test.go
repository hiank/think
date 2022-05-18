package think_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hiank/think"
	"github.com/hiank/think/kube"
	"github.com/hiank/think/store"
	"github.com/hiank/think/store/db"

	"gotest.tools/v3/assert"
)

type testKvStore struct {
	m map[string]any
}

// func (ts *testKvStore)
func (ts *testKvStore) Set(k string, v any) (err error) {
	if k == "" || v == nil {
		return fmt.Errorf("invalid key or value: %s : %v", k, v)
	}
	ts.m[k] = v
	return nil
}

// Get retrieves the value for the given key.
func (ts *testKvStore) Scan(k string, v any) (found bool, err error) {
	if k == "" || v == nil {
		return false, fmt.Errorf("invalid key or value")
	}
	stv, found := ts.m[k]
	if !found {
		return found, fmt.Errorf("cannot found value for given key")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return true, fmt.Errorf("cannot convert value to copy interface")
	}
	mrv := reflect.ValueOf(stv)
	for mrv.Kind() == reflect.Ptr {
		mrv = mrv.Elem()
	}
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	//
	rv.Set(mrv)
	return
}

// Del deletes the stored value for the given key.
func (ts *testKvStore) Del(k string, outs ...any) error {
	if k == "" {
		return fmt.Errorf("invalid key")
	}
	if _, ok := ts.m[k]; !ok {
		return fmt.Errorf("not found")
	}
	delete(ts.m, k)
	return nil
}

func (ts *testKvStore) Close() (err error) {
	return nil
}

var testKvDialer = db.FuncDialer[string](func(c context.Context, do ...db.DialOption) (store.Dictionary[string], error) {
	return &testKvStore{m: map[string]any{}}, nil
})

func TestSetUnique(t *testing.T) {
	t.Run("call Set panic without Awake", func(t *testing.T) {
		defer func(t *testing.T) {
			rcv := recover()
			assert.Assert(t, rcv != nil, "must call Awake before call Set")
		}(t)
		think.Set()
	})
	// t.Run("call Destroy panic without Awake", func(t *testing.T) {
	// 	defer func(t *testing.T) {
	// 		rcv := recover()
	// 		assert.Assert(t, rcv != nil, "must call Awake before call Destroy")
	// 	}(t)
	// 	// think.Destroy()
	// 	think.Set().Close()
	// })
	// t.Run("call Awake panice without necessary options", func(t *testing.T) {
	// 	defer func(t *testing.T) {
	// 		rcv := recover()
	// 		assert.Assert(t, rcv != nil, "must contians necessary options")
	// 	}(t)
	// 	think.Awake()
	// })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	suc := think.Awake(think.WithTODO(ctx), think.WithDB(think.DB{Tag: 1, Dialer: testKvDialer, Opts: []db.DialOption{db.WithAddr("localhost:30211")}}))
	assert.Assert(t, suc)
	suc = think.Awake(think.WithTODO(ctx), think.WithDB(think.DB{Tag: 1, Dialer: testKvDialer, Opts: []db.DialOption{db.WithAddr("localhost:30211")}}))
	assert.Assert(t, !suc, "only the first call works")

	unique := think.Set()
	assert.Assert(t, unique != nil)
	assert.Assert(t, unique == think.Set(), "set is singleston")
	assert.Equal(t, think.Set().Sys(), think.Set().Sys(), "set's part is singleston")
	assert.Equal(t, think.Set().TODO(), think.Set().TODO(), "set's part is singleston")
	assert.Equal(t, think.Set().Nats(), think.Set().Nats(), "set's part is singleston")
	db1, found := think.Set().DB(1)
	assert.Assert(t, found)
	db12, found := think.Set().DB(1)
	assert.Assert(t, found)
	assert.Equal(t, db1, db12, "set's part is singleston")

	think.Set().Close()
	suc = think.Awake(think.WithTODO(ctx), think.WithDB(think.DB{Tag: 1, Dialer: testKvDialer, Opts: []db.DialOption{db.WithAddr("localhost:30211")}}))
	assert.Assert(t, suc)

	assert.Assert(t, unique != think.Set(), "last value destoryed, new value not same as last value")
	assert.Assert(t, unique.Sys() != think.Set().Sys(), "set's part is singleston")
	assert.Assert(t, unique.TODO() != think.Set().TODO(), "set's part is singleston")
	db2, found := think.Set().DB(1)
	assert.Assert(t, found)
	assert.Assert(t, db1 != db2)
	// assert.Assert(t, unique.DBS() != think.Set().DBS(), "set's part is singleston")
	assert.Assert(t, unique.Nats() == think.Set().Nats(), "nats is nil")
	assert.Assert(t, unique.Nats() == nil, "nats is nil")
}

func TestMap(t *testing.T) {
	m := make(map[int]int)
	var i any = m
	_, ok := i.(map[int]any)
	assert.Assert(t, !ok)

	rv := reflect.ValueOf(m)
	assert.Equal(t, rv.Kind(), reflect.Map)

	rt := rv.Type().Elem()
	assert.Equal(t, rt.Kind(), reflect.Int)

	rt = reflect.TypeOf(m).Elem()
	assert.Equal(t, rt.Kind(), reflect.Int)
	// t.Log(rv.Type().Elem())
}

func TestOptions(t *testing.T) {
	dopt := think.Export_defaultOptions()
	assert.Equal(t, dopt.NatsUrl(), kube.NatsUrl())
	assert.Equal(t, dopt.TODO(), context.TODO())
	assert.Equal(t, len(dopt.Mdb()), 0)

	opts := []think.Option{
		think.WithDB(think.DB{Tag: 1, Dialer: testKvDialer, Opts: []db.DialOption{}}),
		think.WithDB(think.DB{Tag: 3, Dialer: testKvDialer, Opts: []db.DialOption{}}),
		think.WithNatsUrl("nats-url"),
		think.WithTODO(context.Background()),
	}

	for _, opt := range opts {
		opt.Apply(dopt.Options())
	}
	assert.Equal(t, dopt.NatsUrl(), "nats-url")
	assert.Equal(t, dopt.TODO(), context.Background())
	assert.Equal(t, len(dopt.Mdb()), 2)
}
