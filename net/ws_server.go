package net

import (
	"errors"
	"github.com/hiank/think/pb"
	// "context"
	"github.com/golang/glog"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
)


//ServeWS ws服务启动
func ServeWS(addrWithPort string) (err error) {

	if netCtx == nil {
		return errors.New("net.Init should be called first")
	}
	k8s.InitClientPool(netCtx, new(k8sReadHandler))
	err = ws.ListenAndServeWS(netCtx, addrWithPort, new(wsReadHandler))
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
func (wh wsReadHandler) Handle(m *pb.Message) (err error) {

	if !k8s.GetClientPool().CheckConnected(m.GetKey(), m.GetToken()) {

		errChan := make(chan error)
		go wh.listen(m.GetKey(), m.GetToken(), errChan)

		var ok bool
		if err, ok = <-errChan; ok {
			glog.Warningln(err)
			return
		}
	}
	k8s.GetClientPool().Post(m)
	return nil
}

func (wh wsReadHandler) listen(key, token string, errChan chan error) {

	cc, err := k8s.DailToCluster(k8s.TypeKubIn, key)
	if err != nil {
		errChan <- err
		return
	}
	defer cc.Close()

	c := pool.NewConnWithDerivedToken(key, token, k8s.NewClientHandler(cc))
	if c == nil {
		errChan <- errors.New("NewConnWithDerivedToken keyed : " + key + "error")
		return
	}
	defer c.GetToken().Cancel()

	k8s.GetClientPool().Push(c)
	close(errChan)
	k8s.GetClientPool().Listen(c)
}
