package mod

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hiank/think"
	"github.com/hiank/think/db"
	"github.com/hiank/think/set"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoCLI = &mongoCLI{conf: new(db.MongoConf), defaultUri: "redis-master"}
)

var defaultMongoConf = `{
	"mongo.TimeoutSecond": 10,
	"mongo.Uri": %s
}`

type mongoCLI struct {
	*mongo.Client
	conf       *db.MongoConf
	defaultUri string //NOTE: 默认url

	think.IgnoreOnDestroy
}

func (cli *mongoCLI) Depend() []think.Module {
	return []think.Module{Config}
}

//OnCreate 此阶段，需要把配置数据注册到ConfigMod
func (cli *mongoCLI) OnCreate(ctx context.Context) error {
	strConf := fmt.Sprintf(defaultMongoConf, cli.defaultUri)
	json.Unmarshal([]byte(strConf), cli.conf)
	Config.SignUpValue(set.JSON, cli.conf)
	return nil
}

func (cli *mongoCLI) OnStart(ctx context.Context) (err error) {
	cli.Client, err = db.NewVerifiedMongoCLI(ctx, cli.conf)
	return err
}

func (cli *mongoCLI) OnStop() {
	cli.Client.Disconnect(context.Background())
}
