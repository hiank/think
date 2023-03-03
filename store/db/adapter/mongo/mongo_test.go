package mongo_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hiank/think/doc"
	"github.com/hiank/think/pbtest"
	"github.com/hiank/think/store"
	"github.com/hiank/think/store/db"
	mgo "github.com/hiank/think/store/db/adapter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gotest.tools/v3/assert"
)

var (
	mongoServerBin, _  = filepath.Abs("testdata/mongod/mongod")
	mongoServerConf, _ = filepath.Abs("testdata/mongod/mongod.conf")
)

func mongoDir(port string) (string, error) {
	dir, err := filepath.Abs(filepath.Join("testdata", "instances", port))
	if err != nil {
		return "", err
	}
	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0775); err != nil {
		return "", err
	}
	os.Mkdir(filepath.Join(dir, "dbpath"), 0775)
	os.Mkdir(filepath.Join(dir, "log"), 0775)
	return dir, nil
}

func startMongo(port string, args ...string) (*os.Process, error) {
	dir, err := mongoDir(port)
	if err != nil {
		return nil, err
	}

	baseArgs := []string{"-f", mongoServerConf, "--port", port, "--dbpath", filepath.Join(dir, "dbpath"), "--logpath", filepath.Join(dir, "log", "mongod.log")}
	return execCmd(mongoServerBin, append(baseArgs, args...)...)
}

func execCmd(name string, args ...string) (*os.Process, error) {
	cmd := exec.Command(name, args...) //append([]string{"-c", "sudo " + name}, args...)...)
	if testing.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Process, cmd.Start()
}

type testMongoStruct struct {
	Name string
	Age  uint
	Id   uint
	Lv   uint
}

type testMongoD struct {
	Id  uint
	Obj *testMongoStruct
}

func funcTestMongoDriver(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:30222"))
	assert.Assert(t, err == nil, err)
	defer cli.Disconnect(ctx)

	kvdb := cli.Database("hi")
	coll := kvdb.Collection("gamer")
	_, err = coll.InsertOne(ctx, bson.D{{Key: "name", Value: "ws"}, {Key: "age", Value: 18}, {Key: "Id", Value: 25}, {Key: "Lv", Value: 11}}) //&testMongoStruct{Name: "ws", Age: 18, Id: 25, Lv: 11})
	assert.Assert(t, err == nil, err)

	var val testMongoStruct
	rlt := coll.FindOne(ctx, bson.D{{Key: "name", Value: "ws"}})
	rlt.Decode(&val)
	assert.Equal(t, val.Name, "ws", val.Name)
	assert.Equal(t, val.Age, uint(18))
	assert.Equal(t, val.Lv, uint(11))
	assert.Equal(t, val.Id, uint(25))

	_, err = coll.InsertOne(ctx, bson.D{{Key: "id", Value: 12}, {Key: "obj", Value: &val}})
	assert.Assert(t, err == nil, err)

	rlt = coll.FindOne(ctx, bson.D{{Key: "id", Value: 12}})
	assert.Assert(t, rlt.Err() == nil, rlt.Err())
	var d testMongoD
	rlt.Decode(&d)
	assert.Equal(t, d.Obj.Name, "ws")
	assert.Equal(t, d.Obj.Age, uint(18))
	assert.Equal(t, d.Obj.Lv, uint(11))
	assert.Equal(t, d.Obj.Id, uint(25))

	_, err = coll.InsertOne(ctx, bson.E{Key: "id", Value: 12})
	assert.Assert(t, err == nil, err)

}

func TestDictionaryPB(t *testing.T) {
	proc, err := startMongo("30222")
	assert.Assert(t, err == nil, err)
	defer proc.Kill()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d, _ := mgo.Dial(ctx, db.WithDB("test"), db.WithAddr("mongodb://localhost:30222"))

	var jk store.Jsonkey //= mgo.DefaultJsonkey()
	var outVal1 pbtest.Test1
	found, err := d.Scan(jk, &outVal1)
	assert.Assert(t, !found)
	assert.Equal(t, err, mgo.ErrNonCollectionOrDocument)

	// (&jk).Encode(store.JsonkeyPair{})
	jk = mgo.DefaultJsonkey()

	err = d.Set(jk, "not proto.Message")
	assert.Assert(t, err != nil, "value to set must be a proto.Message")

	err = d.Set(jk, &pbtest.Test1{Name: "hiank"})
	assert.Assert(t, err == nil, err)

	var outVal2 pbtest.Test2
	found, err = d.Scan(jk, &outVal2)
	assert.Assert(t, found, err)
	assert.Assert(t, err == nil, "protobuf 反序列化时，不同类型也可能会尝试执行，返回的结果不可信")
	assert.Equal(t, outVal2.GetHope(), "")

	found, err = d.Scan(jk, &outVal1)
	assert.Assert(t, found)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, outVal1.GetName(), "hiank")

	err = d.Del(jk)
	assert.Assert(t, err == nil, err)

	found, _ = d.Scan(jk, &outVal1)
	assert.Assert(t, !found)

	err = d.Close()
	assert.Assert(t, err == nil, err)
}

type testDBJson struct {
	Name string `json:"tname"`
	Age  int32  `json:"tage"`
	Lv   int32
	Id   int32
}

type testDBJson2 struct {
	Name string `json:"tname"`
	Age  int    `json:"tage"`
	Lv   int    `json:"tlv"`
	Id   int    `json:"tid"`
}

func funcTestDictionaryJson(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jsonCoder, _ := doc.NewBytesCoder([]byte{}, doc.FormatJson)
	d, _ := mgo.Dial(ctx, db.WithDB("test"), db.WithAddr("mongodb://localhost:30222"), db.WithCoder(jsonCoder))

	var jk store.Jsonkey
	(&jk).Encode(store.JsonkeyPair{K: mgo.JsonkeyCollection, V: "51"}, store.JsonkeyPair{K: mgo.JsonkeyDocument, V: "json"})
	var outVal1 testDBJson2
	found, err := d.Scan(jk, &outVal1)
	assert.Assert(t, !found)
	assert.Assert(t, err != nil)

	err = d.Set(jk, testDBJson{Name: "json", Age: 18, Lv: 22})
	assert.Assert(t, err == nil, err)

	found, err = d.Scan(jk, &outVal1)
	assert.Assert(t, found, err)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, outVal1, testDBJson2{
		Name: "json",
		Age:  18,
		Lv:   0,
		Id:   0,
	})

	var outVal3 testDBJson2
	err = d.Del(jk, &outVal3)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, outVal3.Age, 18)
	found, _ = d.Scan(jk, &outVal1)
	assert.Assert(t, !found)

	err = d.Close()
	assert.Assert(t, err == nil, err)
}

func TestRuns(t *testing.T) {
	proc, err := startMongo("30222")
	assert.Assert(t, err == nil, err)
	defer proc.Kill()

	t.Run("mongo-driver", func(t *testing.T) {
		funcTestMongoDriver(t)
	})

	t.Run("Dictionary-Json", func(t *testing.T) {
		funcTestDictionaryJson(t)
	})
}
