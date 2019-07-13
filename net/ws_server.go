package net

import (
	"github.com/hiank/think/pb"
	"context"
	"github.com/golang/glog"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
)


//ServeWS ws服务启动
func ServeWS(ctx context.Context, addrWithPort string) (err error) {

	k8s.InitClientPool(ctx, new(k8sReadHandler))
	err = ws.ListenAndServeWS(ctx, addrWithPort, new(wsReadHandler))
	return
}


type k8sReadHandler func()
//Handle 处理中grpc 远端读到的Message
func (kh k8sReadHandler) Handle(m *pb.Message) error {

	m.Key = "ws"		//NOTE: 将Message Key 转为 'ws'
	glog.Infoln("k8sReadHandler Handle message :", m)
	ws.GetWSPool().Post(m)
	return nil
}


type wsReadHandler func()
//Handle 处理重ws中读到的数据
func (wh wsReadHandler) Handle(m *pb.Message) error {

	clientPool, key := k8s.GetClientPool(), m.GetKey()
	it := pool.NewDefaultIdentifier(key, m.GetToken())
	if !clientPool.CheckConnected(it) {

		cc, err := k8s.DailToCluster(k8s.TypeKubIn, key)
		if err != nil {

			glog.Infoln("dail to cluster : " + err.Error())
			return err
		}
		c := pool.NewDefaultConn(k8s.NewClientHandler(cc, it))
		clientPool.Push(c)
		wait := make(chan bool)
		go func () {
			close(wait)
			clientPool.Listen(c)
		} ()
		<-wait
	}
	clientPool.Post(m)
	return nil
}
