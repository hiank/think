package mongo

import (
	"context"
	"fmt"

	"github.com/hiank/think/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
)

const (
	docKey    string = "_doc_key"
	docVal    string = "_doc_val"
	defaultDB string = "0"
)

type liteDB struct {
	*mongo.Database
	ctx   context.Context
	coder db.BytesCoder
}

func (ld *liteDB) Get(k string, out interface{}) (found bool, err error) {
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl())
	rlt := coll.FindOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}})
	if err = rlt.Err(); err == nil {
		found, err = true, ld.updateErr(ld.decode(rlt, out))
	}
	return found, ld.updateErr(err)
}

func (ld *liteDB) Set(k string, v interface{}) (err error) {
	bytes, err := ld.coder.Encode(v)
	if err == nil {
		kconv := newKeyConv(k)
		coll := ld.Collection(kconv.GetColl())
		_, err = coll.InsertOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}, {Key: docVal, Value: string(bytes)}})
	}
	return
}

func (ld *liteDB) Del(k string, outs ...interface{}) (err error) {
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl())
	// _, err = coll.DeleteOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}})
	rlt := coll.FindOneAndDelete(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}})
	if err = rlt.Err(); err == nil {
		for _, out := range outs {
			if terr := ld.decode(rlt, out); terr != nil {
				klog.Warning(terr)
			}
		}
	}
	return ld.updateErr(err)
}

func (ld *liteDB) decode(rlt *mongo.SingleResult, out interface{}) (err error) {
	var m bson.M
	if err = rlt.Decode(&m); err == nil {
		if strVal, ok := m[docVal].(string); ok {
			err = ld.coder.Decode([]byte(strVal), out)
		} else {
			err = fmt.Errorf("cached value not a string: %v", m["_val"])
		}
	}
	return
}

func (ld *liteDB) updateErr(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		err = db.ErrNotFound
	}
	return err
}

func (ld *liteDB) Close() error {
	return ld.Client().Disconnect(ld.ctx)
}
