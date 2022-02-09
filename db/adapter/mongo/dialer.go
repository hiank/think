package mongo

import (
	"context"

	"github.com/hiank/think/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Dialer db.KvDialer = dialer(0)
)

type dialer byte

//Dial connect to mongodb and return connected client or error
func (d dialer) Dial(ctx context.Context, opts ...db.DialOption) (kv db.KvDB, err error) {
	dopts := db.DialOptions(opts...)
	mopt := options.Client().ApplyURI(dopts.Addr).SetConnectTimeout(dopts.DialTimeout)
	if dopts.Account != "" || dopts.Password != "" {
		mopt.SetAuth(options.Credential{
			Username: dopts.Account,
			Password: dopts.Password,
		})
	}
	cli, err := mongo.Connect(ctx, mopt)
	if err == nil {
		kv = &liteDB{
			ctx:      ctx,
			Database: cli.Database(dopts.DB),
		}
	}
	return
}
