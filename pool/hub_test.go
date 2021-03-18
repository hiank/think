package pool_test

import (
	"testing"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"gotest.tools/v3/assert"
)

func TestLimitMux(t *testing.T) {

	lm := &pool.LimitMux{Max: 1}
	assert.Assert(t, !lm.Locked(), "未到限制数时，不锁定")
	assert.Assert(t, lm.Retain(), "未到限制数时，可请求")
	assert.Assert(t, lm.Locked(), "到限制数，锁定")
	assert.Assert(t, !lm.Retain(), "到限制数，不能请求")
}

func TestListMux(t *testing.T) {

	listMux := pool.NewListMux()
	assert.Equal(t, listMux.Shift(), nil, "没有数据时，Shift返回nil")
	listMux.Push(&pb.Message{Key: "test1"})
	listMux.Push(&pb.Message{Key: "test2"})
	assert.Equal(t, listMux.Shift().(*pb.Message).GetKey(), "test1", "有数据时，Shift返回最前面数据")
	assert.Equal(t, listMux.Shift().(*pb.Message).GetKey(), "test2", "")
	assert.Equal(t, listMux.Shift(), nil, "")
}
