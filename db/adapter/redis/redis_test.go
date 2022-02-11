package redis_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/db"
	rdbc "github.com/hiank/think/db/adapter/redis"
	"github.com/hiank/think/doc/testdata"
	"gotest.tools/v3/assert"
)

var (
	redisServerBin, _  = filepath.Abs("testdata/redis/redis-server")
	redisServerConf, _ = filepath.Abs("testdata/redis/redis.conf")
)

func redisDir(port string) (string, error) {
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
	return dir, nil
}

func startRedis(port string, args ...string) (*os.Process, error) {
	dir, err := redisDir(port)
	if err != nil {
		return nil, err
	}

	baseArgs := []string{redisServerConf, "--port", port, "--dir", dir}
	return execCmd(redisServerBin, append(baseArgs, args...)...)
}

func execCmd(name string, args ...string) (*os.Process, error) {
	cmd := exec.Command(name, args...)
	if testing.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Process, cmd.Start()
}

type testDBStruct struct {
	Name string `redis:"name"`
	Age  uint   `redis:"age"`
	Lv   uint   `redis:"lv"`
	Id   uint   `redis:"id"`
}

func TestRedisCli(t *testing.T) {
	proc, err := startRedis("30211")
	if err != nil {
		t.Error(err)
		return
	}
	defer proc.Kill()
	<-time.After(time.Second)

	ctx, cancel := context.WithCancel(context.Background()) //context.Background()
	defer cancel()

	t.Run("redis client", func(t *testing.T) {
		rdbCli := redis.NewClient(&redis.Options{
			DB:       1,
			Password: "",
			Addr:     "localhost:30211",
		})

		_, err := rdbCli.HSet(ctx, "hash", "key1", "hello1", "key2", "hello2").Result()
		assert.Assert(t, err == nil, err)
	})

	t.Run("connect error", func(t *testing.T) {
		// cli, err := rdbc.Dial(ctx, &redis.Options{
		// 	DB:          0,
		// 	Password:    "",
		// 	Addr:        "localhost:30001",
		// 	DialTimeout: time.Second,
		// })
		cli, err := rdbc.Dialer.Dial(ctx, db.WithAddr("localhost:30001"), db.WithDB("0"), db.WithDailTimeout(time.Second))
		assert.Assert(t, cli == nil)
		assert.Assert(t, err != nil)
	})

	t.Run("CRUD-PB", func(t *testing.T) {
		// cli, _ := rdbc.Dial(ctx, &redis.Options{
		// 	DB:       0,
		// 	Password: "",
		// 	Addr:     "localhost:30211",
		// })
		cli, _ := rdbc.Dialer.Dial(ctx, db.WithAddr("localhost:30211"), db.WithDB("0"))
		defer cli.Close()
		err := cli.Set("hs", "key1")
		assert.Assert(t, err != nil, "value must be proto.Message")

		var val1 testdata.Test1
		val1.Name = "val1"
		err = cli.Set("hs", &val1)
		assert.Assert(t, err != nil, "must use PB JSON GOB struct value")
		err = cli.Set("hs", db.PB{V: &val1})
		assert.Assert(t, err == nil, err)

		var outVal1 testdata.Test1
		found, err := cli.Get("key1", db.PB{V: &outVal1})
		assert.Assert(t, !found)
		assert.Equal(t, err, db.ErrNotFound)

		found, err = cli.Get("hs", db.PB{V: &outVal1})
		assert.Assert(t, found)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, outVal1.GetName(), "val1")

		err = cli.Set("hs", db.PB{V: &testdata.Test2{Age: 18}})
		assert.Assert(t, err == nil, err)

		var outVal2 testdata.Test2
		found, err = cli.Get("hs", db.PB{V: &outVal2})
		assert.Assert(t, found)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, outVal2.Age, int32(18))

		err = cli.Del("key1")
		// assert.Assert(t, err == nil, err)
		assert.Equal(t, err, db.ErrNotFound, "delete not existed key")
		found, _ = cli.Get("hs", db.PB{V: &outVal2})
		assert.Assert(t, found)

		var outVal3 testdata.Test2
		err = cli.Del("hs", db.PB{V: &outVal3})
		assert.Assert(t, err == nil, err)
		assert.Equal(t, outVal3.GetAge(), int32(18))
		found, _ = cli.Get("hs", db.PB{V: &outVal2})
		assert.Assert(t, !found)

		err = cli.Close()
		assert.Assert(t, err == nil, err)
	})

	t.Run("CRUD-Json", func(t *testing.T) {
		// cli, _ := rdbc.Dial(ctx, &redis.Options{
		// 	DB:       0,
		// 	Password: "",
		// 	Addr:     "localhost:30211",
		// })
		cli, _ := rdbc.Dialer.Dial(ctx, db.WithAddr("localhost:30211"), db.WithDB("0"))
		defer cli.Close()
		// err := cli.Set("hs", "key1")
		// assert.Assert(t, err == nil, err)

		// var val1 testDBStruct//testdata.Test1
		// val1.Name = "val1"
		val1 := &testDBStruct{
			Name: "hiank",
			Age:  19,
			Lv:   25,
			Id:   11,
		}
		err = cli.Set("hs", &val1)
		assert.Assert(t, err != nil, err)
		err = cli.Set("hs", db.JSON{V: &val1})
		assert.Assert(t, err == nil, err)

		var outVal1 testDBStruct
		found, err := cli.Get("key1", db.JSON{V: &outVal1})
		assert.Assert(t, !found)
		assert.Assert(t, err != nil)

		found, err = cli.Get("hs", db.JSON{V: &outVal1})
		assert.Assert(t, found)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, outVal1, *val1)

		err = cli.Del("hs")
		assert.Assert(t, err == nil)
		found, err = cli.Get("hs", db.JSON{V: &outVal1})
		assert.Assert(t, !found)
		assert.Assert(t, err != nil, err)

		err = cli.Close()
		assert.Assert(t, err == nil, err)
	})
}
