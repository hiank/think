package mongo

import (
	"context"
	"fmt"

	"github.com/hiank/think/data/db"
	"github.com/hiank/think/doc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
	// mopts ""
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
	// opt     *Options
	opts     *options
	docMaker doc.BytesMaker
}

func (ld *liteDB) Get(k string, v interface{}) (found bool, err error) {
	var m bson.M
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl(), ld.opts.collectionOpts...)
	rlt := coll.FindOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}}, ld.opts.findOneOpts...)
	if err = rlt.Err(); err != nil {
		return
	}
	found = true
	if err = rlt.Decode(&m); err == nil {
		if strVal, ok := m[docVal].(string); ok {
			err = ld.docMaker.Make([]byte(strVal)).Decode(v)
		} else {
			err = fmt.Errorf("cached value not a string: %v", m["_val"])
		}
	}
	return
}

func (ld *liteDB) Set(k string, v interface{}) (err error) {
	doc := ld.docMaker.Make(nil)
	if err = doc.Encode(v); err == nil {
		kconv := newKeyConv(k)
		coll := ld.Collection(kconv.GetColl(), ld.opts.collectionOpts...)
		_, err = coll.InsertOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}, {Key: docVal, Value: doc.Val()}}, ld.opts.insertOneOpts...)
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

func NewKvDB(ctx context.Context, docMaker doc.BytesMaker, opts ...Option) db.KvDB {
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
		docMaker: docMaker,
	}
}
