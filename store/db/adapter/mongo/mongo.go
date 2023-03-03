package mongo

import (
	"context"
	"io"

	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
	"github.com/hiank/think/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
)

const (
	docKey string = "_doc_key"
	docVal string = "_doc_val"

	JsonkeyCollection string = "collection"
	JsonkeyDocument   string = "document"

	ErrNonCollectionOrDocument = run.Err("mongo: non collection or document key")
	ErrNotString               = run.Err("mongo: not a string value")
	ErrNotFound                = run.Err("mongo: not found")
)

func DefaultJsonkey() (jk store.Jsonkey) {
	(&jk).Encode(store.JsonkeyPair{K: JsonkeyCollection, V: "0"}, store.JsonkeyPair{K: JsonkeyDocument, V: "default"})
	return
}

// type StoreKey store.Jsonkey

type liteDB struct {
	db    *mongo.Database
	ctx   context.Context
	coder doc.Coder
	io.Closer
}

func (ld *liteDB) Scan(k store.Jsonkey, out any) (found bool, err error) {
	ck, dk, foundKey, err := ld.checkey(k)
	if foundKey {
		coll := ld.db.Collection(ck)
		rlt := coll.FindOne(ld.ctx, bson.D{{Key: docKey, Value: dk}})
		if err = rlt.Err(); err == nil {
			found, err = true, ld.updateErr(ld.decode(rlt, out))
		}
	}
	return
}

func (ld *liteDB) Set(k store.Jsonkey, v any) (err error) {
	if err = ld.coder.Encode(v); err != nil {
		return
	}
	ck, dk, found, err := ld.checkey(k)
	if found {
		coll := ld.db.Collection(ck)
		_, err = coll.InsertOne(ld.ctx, bson.D{{Key: docKey, Value: dk}, {Key: docVal, Value: string(ld.coder.Bytes())}})
	}
	return ld.updateErr(err)
}

func (ld *liteDB) Del(k store.Jsonkey, outs ...any) (err error) {
	ck, dk, found, err := ld.checkey(k)
	if found {
		coll := ld.db.Collection(ck)
		rlt := coll.FindOneAndDelete(ld.ctx, bson.D{{Key: docKey, Value: dk}})
		if err = rlt.Err(); err == nil {
			for _, out := range outs {
				if terr := ld.decode(rlt, out); terr != nil {
					klog.Warning(terr)
				}
			}
		}
	}
	return ld.updateErr(err)
}

func (ld *liteDB) decode(rlt *mongo.SingleResult, out any) (err error) {
	var m bson.M
	if err = rlt.Decode(&m); err == nil {
		if strVal, ok := m[docVal].(string); ok {
			ld.coder.Encode([]byte(strVal))
			err = ld.coder.Decode(out)
		} else {
			err = ErrNotString //fmt.Errorf("cached value not a string: %v", m["_val"])
		}
	}
	return
}

func (ld *liteDB) checkey(k store.Jsonkey) (ck, dk string, ok bool, err error) {
	if ck, ok = k.Get(JsonkeyCollection); ok {
		if dk, ok = k.Get(JsonkeyDocument); ok {
			return
		}
	}
	return "", "", false, ErrNonCollectionOrDocument
}

func (ld *liteDB) updateErr(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		err = ErrNotFound
	}
	return err
}

// func (ld *liteDB) Close() error {
// 	return ld.Client().Disconnect(ld.ctx)
// }
