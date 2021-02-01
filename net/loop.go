package net

import (
	"context"

	"github.com/hiank/think/pool"
	"k8s.io/klog/v2"
)

func loopAccept(ctx context.Context, accepter Accepter, handler func(Conn)) {

	for {
		val, err := accepter.Accept()
		select {
		case <-ctx.Done():
			return
		default:
		}
		switch err {
		case nil:
			if val != nil {
				handler(val)
			} else {
				klog.Warning("conn by Accept is nil") //NOTE: 收消息错误，直接退出，需要验证下，发一个空消息会怎样
			}
		default:
			klog.Warning(err.Error()) //NOTE: 收消息错误，直接退出，需要验证下，发一个空消息会怎样
			return
		}
	}
}

func loopRecv(ctx context.Context, reciver Reciver, handler pool.Handler) {

	for {
		val, err := reciver.Recv()
		select {
		case <-ctx.Done():
			return
		default:
		}
		switch err {
		case nil:
			handler.Handle(val)
		default:
			klog.Warning(err.Error()) //NOTE: 收消息错误，直接退出，需要验证下，发一个空消息会怎样
			return
		}
	}
}

// //Getter 获取资源
// type Getter interface {
// 	Get() (interface{}, error)
// }

// func loopWork(ctx context.Context, getter Getter, handler pool.Handler, closers ...io.Closer) {

// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	go func() {
// 		<-ctx.Done()
// 		for _, closer := range closers {
// 			closer.Close()
// 		}
// 	}()

// 	for {
// 		val, err := getter.Get()
// 		switch err {
// 		case io.EOF:
// 			return
// 		case nil:
// 			handler.Handle(val)
// 		default:
// 			klog.Warning(err.Error())
// 		}
// 	}
// }

// type recvGetter func() (*pb.Message, error)

// func (rg recvGetter) Get() (interface{}, error) {
// 	return rg()
// }

// type connGetter func() (Conn, error)

// func (cg connGetter) Get() (interface{}, error) {
// 	return cg()
// }
