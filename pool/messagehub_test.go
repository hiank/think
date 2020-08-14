package pool_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
)

// //test api Push
// func Test(t *testing.T) {

// }

type intTypeHandler int

func (iyh intTypeHandler) Handle(msg *pool.Message) error {
	// return fyh(msg)
	// <-time.After(time.Millisecond * 100)
	time.Sleep(time.Millisecond * 100)
	// fmt.Print(msg.ToString(), ",")
	return nil
}

//test api PushWithBack
func TestPushWithBack(t *testing.T) {

}

//test api LockReq
func TestDoActive(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	hub := pool.NewMessageHub(ctx, intTypeHandler(1))
	for i:=0; i<1000; i++{
		hub.Push(pool.NewMessage(nil, token.GetBuilder().Get(strconv.Itoa(i))))
	}
	<-time.After(1)
	fmt.Println("waitted 1s")

	hub.DoActive()
	<-time.After(time.Millisecond * 110)
	cancel()
	<-time.After(time.Second)
	fmt.Println("waitted next 2s")
}

//test messagehub working state
func TestMessageHubWorking(t *testing.T) {

}
