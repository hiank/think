package mongo

import (
	"context"
	"fmt"

	"github.com/hiank/think/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	docKey    string = "_doc_key"
	docVal    string = "_doc_val"
	defaultDB string = "0"
)

func defaultOptions() options {
	return options{
		dbName:         defaultDB,
		clientOpts:     []*mopts.ClientOptions{},
		databaseOpts:   []*mopts.DatabaseOptions{},
		collectionOpts: []*mopts.CollectionOptions{},
		findOneOpts:    []*mopts.FindOneOptions{},
		deleteOpts:     []*mopts.DeleteOptions{},
	}
}

type liteDB struct {
	ctx context.Context
	*mongo.Database
	opts  *options
	coder db.BytesCoder
}

func (ld *liteDB) Get(k string, out interface{}) (found bool, err error) {
	var m bson.M
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl(), ld.opts.collectionOpts...)
	rlt := coll.FindOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}}, ld.opts.findOneOpts...)
	if rlt.Err() != nil {
		return false, rlt.Err()
	}
	if err = rlt.Decode(&m); err == nil {
		if strVal, ok := m[docVal].(string); ok {
			err = ld.coder.Decode([]byte(strVal), out)
		} else {
			err = fmt.Errorf("cached value not a string: %v", m["_val"])
		}
	}
	return true, err
}

func (ld *liteDB) Set(k string, v interface{}) (err error) {
	bytes, err := ld.coder.Encode(v)
	if err == nil {
		kconv := newKeyConv(k)
		coll := ld.Collection(kconv.GetColl(), ld.opts.collectionOpts...)
		_, err = coll.InsertOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}, {Key: docVal, Value: string(bytes)}}, ld.opts.insertOneOpts...)
	}
	return
}

func (ld *liteDB) Delete(k string) (err error) {
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl(), ld.opts.collectionOpts...)
	_, err = coll.DeleteOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}}, ld.opts.deleteOpts...)
	return
}

func (ld *liteDB) Close() error {
	return ld.Client().Disconnect(ld.ctx)
}

func NewKvDB(ctx context.Context, opts ...Option) db.KvDB {
	dopts := defaultOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	cli, err := mongo.Connect(ctx, dopts.clientOpts...)
	if err != nil {
		panic(err)
	}
	return &liteDB{
		ctx:      ctx,
		Database: cli.Database(dopts.dbName, dopts.databaseOpts...),
		opts:     &dopts,
	}
}
