package mongo

import (
	"context"

	"github.com/hiank/think/run"
	"github.com/hiank/think/store"
	"github.com/hiank/think/store/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog/v2"
)

//Dial connect to mongodb and return connected client or error
func Dial(ctx context.Context, opts ...db.DialOption) (d store.Dictionary[store.Jsonkey], err error) {
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
		mdb := cli.Database(dopts.DB)
		ctx, closer := run.StartHealthyMonitoring(ctx, func() {
			klog.Warning(mdb.Client().Disconnect(ctx))
		})
		d = &liteDB{
			ctx:    ctx,
			db:     mdb,
			coder:  dopts.Coder,
			Closer: closer,
		}
	}
	return
}
