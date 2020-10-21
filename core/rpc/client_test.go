package rpc_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/hiank/think/core/rpc"
	"github.com/hiank/think/settings"
	"gotest.tools/assert"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	td "github.com/hiank/think/core/rpc/testdata"
)

type testClientReadHandler struct {
	handleNum int
}

func (tr *testClientReadHandler) HandleGet(msg *pb.Message) (*pb.Message, error) {

	fmt.Println("HandleGet")
	return msg, nil
}

func (tr *testClientReadHandler) HandlePost(msg *pb.Message) (err error) {
	fmt.Println("HandlePost")

	val := &td.P_Example{}
	if err = ptypes.UnmarshalAny(msg.GetValue(), val); err != nil {
		return
	}
	if val.GetValue() != "post" {
		err = errors.New(val.GetValue())
	}
	fmt.Println("HandlePost:", err)
	return
}

func (tr *testClientReadHandler) Handle(msg core.Message) error {

	// tr.handleNum++
	// fmt.Println("handle num : ", tr.handleNum)

	var writer rpc.Writer
	return writer.Handle(msg) //NOTE: 将收到的消息发回去
}

func TestClient(t *testing.T) {

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	wait := new(sync.WaitGroup)

	wait.Add(1)

	go startOneServer(ctx, wait, new(testClientReadHandler))

	client := rpc.NewClient(ctx, "test")
	cc, err := client.Dial(fmt.Sprintf("localhost:%v", settings.GetSys().K8sPort))

	assert.Assert(t, err, nil, "必须能成功连上服务端")

	defer cc.Close()

	send := func(msg proto.Message) error {
		// t.Log("in send:", msg)
		a, _ := ptypes.MarshalAny(msg)

		return client.Send(&pb.Message{Key: "token", Value: a})
	}

	t.Run("TypeGET", func(t *testing.T) {

		gdata := &td.G_Example{Value: "g"}
		assert.Assert(t, send(gdata), nil, "必须能成功发送")
		msg, _ := client.Recv()
		// assert.Equal(t, fmt.Sprintf("%v", msg), fmt.Sprintf("%v", gdata), "发送的消息与返回的消息值需要是一样")
		assert.Assert(t, fmt.Sprintf("%p", msg) != fmt.Sprintf("%p", gdata), "发送的消息与返回的消息是两个不同的对象，数据地址需要不同")
	})

	t.Run("TypePOST", func(t *testing.T) {

		assert.Assert(t, send(&td.P_Example{Value: "post"}) == nil, "需要成功发送消息")
		assert.Assert(t, strings.Contains(send(&td.P_Example{Value: "errorTest"}).Error(), "errorTest"), "发送错误消息，期望返回错误")
	})

	t.Run("TypeLink", func(t *testing.T) {

		arr := [10000]int{}
		num := len(arr)
		go func() {
			notice, noticeBefore := make(chan bool), make(chan bool)
			for i := 0; i < num; i++ {
				go func(i int) {
					noticeBefore <- true
					send(&td.S_Example{Value: strconv.Itoa(i)})
					notice <- true
				}(i)
			}
			var sendCnt, beforeCnt int
		L:
			for {
				select {
				case <-noticeBefore:
					beforeCnt++
				case <-notice:
					sendCnt++
				case <-time.After(time.Second * 18):
					break L
				}
			}
			t.Log("sendCnt", sendCnt, beforeCnt)
		}()

		var recvCnt int
		for {
			recvMsg, err := client.Recv()
			if err != nil {
				break
			}
			val := new(td.S_Example)
			err = ptypes.UnmarshalAny(recvMsg.GetValue(), val)

			assert.Assert(t, err == nil, "需要能正确收到消息", err)
			i, _ := strconv.Atoi(val.GetValue())
			arr[i] = -1
			recvCnt++
			if num == recvCnt {
				break
			}
		}
		for i, val := range arr {
			if val != -1 {
				t.Log("idx not recv:", i)
			}
			// assert.Equal(t, val, 0, "需要所有的位都置0")
		}
	})
	cancel()
	wait.Wait()
}

func TestClientClose(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	// testLink(t, "")

	waitServe := new(sync.WaitGroup)
	waitServe.Add(1)
	go startOneServer(ctx, waitServe, new(testClientReadHandler))

	client := rpc.NewClient(ctx, "test")
	_, err := client.Dial(fmt.Sprintf("localhost:%v", settings.GetSys().K8sPort))
	assert.Assert(t, err, nil, "必须能成功连上服务端")

	a, _ := ptypes.MarshalAny(&td.S_Example{Value: "moreKey"})
	err = client.Send(&pb.Message{Key: "token", Value: a})
	assert.Assert(t, err, nil, "需要能成功发送消息", err)

	recvMsg, err := client.Recv()
	assert.Assert(t, err, nil, "需要能成功收到消息", err)

	val := new(td.S_Example)
	err = ptypes.UnmarshalAny(recvMsg.GetValue(), val)
	assert.Assert(t, err == nil, "需要能正确收到消息", err)
	assert.Equal(t, "moreKey", val.GetValue())

	client.Close()

	time.Sleep(time.Second)

	cancel()
	waitServe.Wait()
}

func TestMultiClientContainSomeToken(t *testing.T) {

	wait, maxNum := new(sync.WaitGroup), 200
	wait.Add(maxNum)

	ctx, cancel := context.WithCancel(context.Background())
	var testLink = func(t *testing.T, moreKey string) {

		client := rpc.NewClient(ctx, "test")
		_, err := client.Dial(fmt.Sprintf("localhost:%v", settings.GetSys().K8sPort))
		assert.Assert(t, err, nil, "必须能成功连上服务端")

		a, _ := ptypes.MarshalAny(&td.S_Example{Value: moreKey})
		err = client.Send(&pb.Message{Key: "token", Value: a})
		assert.Assert(t, err, nil, "需要能成功发送消息", err)

		recvMsg, err := client.Recv()
		assert.Assert(t, err, nil, "需要能成功收到消息", err)

		val := new(td.S_Example)
		err = ptypes.UnmarshalAny(recvMsg.GetValue(), val)
		assert.Assert(t, err == nil, "需要能正确收到消息", err)
		assert.Equal(t, moreKey, val.GetValue())

		client.Close()

		time.Sleep(time.Second)
		wait.Done()
	}
	// testLink(t, "")

	waitServe := new(sync.WaitGroup)
	waitServe.Add(1)
	go startOneServer(ctx, waitServe, new(testClientReadHandler))

	for i := 0; i < maxNum; i++ {
		go testLink(t, strconv.Itoa(i))
	}

	// time.Sleep(time.Second * 5)
	wait.Wait()

	cancel()
	waitServe.Wait()
}
