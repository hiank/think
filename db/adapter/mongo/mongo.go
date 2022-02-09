package mongo

import (
	"context"
	"fmt"

	"github.com/hiank/think/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	var m bson.M
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl())
	rlt := coll.FindOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}})
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
		coll := ld.Collection(kconv.GetColl())
		_, err = coll.InsertOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}, {Key: docVal, Value: string(bytes)}})
	}
	return
}

func (ld *liteDB) Delete(k string) (err error) {
	kconv := newKeyConv(k)
	coll := ld.Collection(kconv.GetColl())
	_, err = coll.DeleteOne(ld.ctx, bson.D{{Key: docKey, Value: kconv.GetDoc()}})
	return
}

func (ld *liteDB) Close() error {
	return ld.Client().Disconnect(ld.ctx)
}
