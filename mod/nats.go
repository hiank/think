package mod

import (
	"context"
	"encoding/json"

	"github.com/hiank/think"
	"github.com/hiank/think/mod/modex"
	"github.com/hiank/think/net/mq"
	"github.com/hiank/think/set"
	"github.com/nats-io/nats.go"
)

var (
	NatsCLI = &natsCLI{addrGetter: new(NatsAddr)}
)

var defaultNatsConf = `{
	"nats.Addr": {
		"Addr": "nats"
	}
}`

type NatsAddr struct {
	*modex.Addr `json:"nats.Addr"`
}

func (natsAddr *NatsAddr) Get() *modex.Addr {
	return natsAddr.Addr
}

type natsCLI struct {
	*nats.Conn
	addrGetter modex.AddrGetter

	think.IgnoreOnDestroy
}

func (cli *natsCLI) Depend() []think.Module {
	return []think.Module{KubesetIn, Config}
}

//OnCreate 此阶段，需要把配置数据注册到ConfigMod
func (cli *natsCLI) OnCreate(ctx context.Context) error {
	json.Unmarshal([]byte(defaultNatsConf), cli.addrGetter)
	Config.SignUpValue(set.JSON, cli.addrGetter)
	return nil
}

func (cli *natsCLI) OnStart(ctx context.Context) (err error) {
	addr, err := modex.ParseAddr(cli.addrGetter.Get().Value, KubesetIn)
	if err == nil {
		cli.Conn, err = mq.NewNatsConn(addr)
	}
	return
}

func (cli *natsCLI) OnStop() {
	cli.Close()
}
