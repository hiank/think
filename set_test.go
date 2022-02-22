package think_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hiank/think"
	"github.com/hiank/think/db"
	"gotest.tools/v3/assert"
)

type testKvStore struct {
	m map[string]interface{}
}

// func (ts *testKvStore)
func (ts *testKvStore) Set(k string, v interface{}) (err error) {
	if k == "" || v == nil {
		return fmt.Errorf("invalid key or value: %s : %v", k, v)
	}
	ts.m[k] = v
	return nil
}

// Get retrieves the value for the given key.
func (ts *testKvStore) Get(k string, v interface{}) (found bool, err error) {
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
func (ts *testKvStore) Del(k string, outs ...interface{}) error {
	if k == "" {
		return fmt.Errorf("invalid key")
	}
	if _, ok := ts.m[k]; !ok {
		return db.ErrNotFound
	}
	delete(ts.m, k)
	return nil
}

func (ts *testKvStore) Close() (err error) {
	return nil
}

var testKvDialer = think.FuncKvDialer(func(c context.Context, do ...db.DialOption) (db.KvDB, error) {
	return &testKvStore{m: map[string]interface{}{}}, nil
})

func TestSetUnique(t *testing.T) {
	t.Run("call Set panic without Awake", func(t *testing.T) {
		defer func(t *testing.T) {
			rcv := recover()
			assert.Assert(t, rcv != nil, "must call Awake before call Set")
		}(t)
		think.Set()
	})
	t.Run("call Destroy panic without Awake", func(t *testing.T) {
		defer func(t *testing.T) {
			rcv := recover()
			assert.Assert(t, rcv != nil, "must call Awake before call Destroy")
		}(t)
		think.Destroy()
	})
	// t.Run("call Awake panice without necessary options", func(t *testing.T) {
	// 	defer func(t *testing.T) {
	// 		rcv := recover()
	// 		assert.Assert(t, rcv != nil, "must contians necessary options")
	// 	}(t)
	// 	think.Awake()
	// })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	suc := think.Awake(think.WithTODO(ctx), think.WithMemKvDialer(testKvDialer, db.WithAddr("localhost:30211")))
	assert.Assert(t, suc)
	suc = think.Awake(think.WithTODO(ctx), think.WithMemKvDialer(testKvDialer, db.WithAddr("localhost:30211")))
	assert.Assert(t, !suc, "only the first call works")

	unique := think.Set()
	assert.Assert(t, unique != nil)
	assert.Assert(t, unique == think.Set(), "set is singleston")
	assert.Equal(t, think.Set().Decoder(), think.Set().Decoder(), "set's part is singleston")
	assert.Equal(t, think.Set().TODO(), think.Set().TODO(), "set's part is singleston")
	assert.Equal(t, think.Set().Nats(), think.Set().Nats(), "set's part is singleston")
	assert.Equal(t, think.Set().DBS(), think.Set().DBS(), "set's part is singleston")

	think.Destroy()
	suc = think.Awake(think.WithTODO(ctx), think.WithMemKvDialer(testKvDialer, db.WithAddr("localhost:30211")))
	assert.Assert(t, suc)

	assert.Assert(t, unique != think.Set(), "last value destoryed, new value not same as last value")
	assert.Assert(t, unique.Decoder() != think.Set().Decoder(), "set's part is singleston")
	assert.Assert(t, unique.TODO() != think.Set().TODO(), "set's part is singleston")
	assert.Assert(t, unique.DBS() != think.Set().DBS(), "set's part is singleston")
	assert.Assert(t, unique.Nats() == think.Set().Nats(), "nats is nil")
	assert.Assert(t, unique.Nats() == nil, "nats is nil")
}

func TestMap(t *testing.T) {
	m := make(map[int]int)
	var i interface{} = m
	_, ok := i.(map[int]interface{})
	assert.Assert(t, !ok)

	rv := reflect.ValueOf(m)
	assert.Equal(t, rv.Kind(), reflect.Map)

	rt := rv.Type().Elem()
	assert.Equal(t, rt.Kind(), reflect.Int)

	rt = reflect.TypeOf(m).Elem()
	assert.Equal(t, rt.Kind(), reflect.Int)
	// t.Log(rv.Type().Elem())
}
