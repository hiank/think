package data_test

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hiank/think/data"
	"github.com/hiank/think/data/db"
	"gotest.tools/v3/assert"
)

// var (
// 	redisServerBin, _  = filepath.Abs("db/testdata/redis/src/redis-server")
// 	redisServerConf, _ = filepath.Abs("db/testdata/redis/redis.conf")
// )

// func redisDir(port string) (string, error) {
// 	dir, err := filepath.Abs(filepath.Join("db/testdata", "instances", port))
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := os.RemoveAll(dir); err != nil {
// 		return "", err
// 	}
// 	if err := os.MkdirAll(dir, 0775); err != nil {
// 		return "", err
// 	}
// 	return dir, nil
// }

// func startRedis(port string, args ...string) (*os.Process, error) {
// 	dir, err := redisDir(port)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = exec.Command("cp", "-f", redisServerConf, dir).Run(); err != nil {
// 		return nil, err
// 	}

// 	baseArgs := []string{filepath.Join(dir, "redis.conf"), "--port", port, "--dir", dir}
// 	return execCmd(redisServerBin, append(baseArgs, args...)...)
// }

// func execCmd(name string, args ...string) (*os.Process, error) {
// 	cmd := exec.Command(name, args...)
// 	if testing.Verbose() {
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 	}
// 	return cmd.Process, cmd.Start()
// }

type testHValue struct {
	IntVal  int    `redis:"intval"`
	BoolVal bool   `redis:"boolval"`
	StrVal  string `redis:"strval"`
	UintVal uint   `redis:"uintval"`
}

func TestStringSetValue(t *testing.T) {
	var str1 string = "name"
	str2 := &str1
	*str2 = "test"

	assert.Equal(t, str1, "test")
	assert.Equal(t, *str2, "test")
}

// type

// func TestBaseRedis(t *testing.T) {
// 	proc, err := startRedis("30211")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer proc.Kill()

// 	ctx := context.Background()
// 	cli := redis.NewClient(&redis.Options{
// 		DB:       0,
// 		Password: "",
// 		Addr:     "localhost:30211",
// 	}) //db.NewVerifiedRedisCLI(ctx, redisConf)

// 	t.Run("lite Get-Set", func(t *testing.T) {
// 		cli.Set(ctx, "testInt", 1, 0)
// 		val, err := cli.Get(ctx, "testInt").Result()

// 		assert.Assert(t, err == nil, err)
// 		assert.Equal(t, val, "1")
// 	})

// 	t.Run("HGet-HSet", func(t *testing.T) {
// 		cli.HMSet(ctx, "h-key", "intval", "-11", "uintval", "119", "boolval", "true")
// 		cmd := cli.HGet(ctx, "h-key", "intval")
// 		assert.Assert(t, cmd.Err() == nil)
// 		rlt, _ := cmd.Result()
// 		assert.Equal(t, rlt, "-11")

// 		allCmd := cli.HGetAll(ctx, "h-key")
// 		m, err := allCmd.Result()
// 		assert.Assert(t, err == nil)
// 		assert.Equal(t, m["intval"], "-11")
// 		assert.Equal(t, m["uintval"], "119")

// 		val := &testHValue{}
// 		err = allCmd.Scan(val)
// 		assert.Assert(t, err == nil)

// 		assert.Equal(t, val.IntVal, -11)
// 		assert.Equal(t, val.UintVal, uint(119))
// 		assert.Equal(t, val.BoolVal, true)
// 		assert.Equal(t, val.StrVal, "")
// 	})
// }

// type testMemoryDB struct {
// }

// func (tm *testMemoryDB) HGet(hkey, fkey string) (db.IParser, error) {
// 	return nil, nil
// }

// func (tm *testMemoryDB) HSet(hkey string, vals ...interface{}) error {
// 	return nil
// }

// func (tm *testMemoryDB) Close() error {
// 	return nil
// }

func TestPushError(t *testing.T) {
	err := data.Export_pushErr(nil, nil)
	assert.Assert(t, err == nil, err)

	err = data.Export_pushErr(nil, errors.New("err1"))
	assert.Equal(t, err.Error(), "err1")

	err = data.Export_pushErr(err, nil)
	assert.Equal(t, err.Error(), "err1")

	err = data.Export_pushErr(err, errors.New("err2"))
	assert.Equal(t, err.Error(), "err1&&err2")
}

func TestRobustDB(t *testing.T) {
	tks := &testKvStore{m: map[string]interface{}{}}
	rdb := data.Export_newRobustDB(tks)
	//
	err := rdb.Set("key", 1)
	assert.Assert(t, err != nil, "invalid key")
	err = rdb.Set(data.KTMem.Encode("key"), 1)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, tks.m["key"], 1)
	// assert.Equal(t, tks.m["key"], )
	var val int
	rdb.Get(data.KTDisk.Encode("key"), &val)
	assert.Equal(t, val, 1)

	err = rdb.Delete(data.KTMix.Encode("key"))
	assert.Assert(t, err == nil, err)
	found, err := rdb.Get(data.KTDisk.Encode("key"), &val)
	assert.Assert(t, !found)
	assert.Assert(t, err != nil)
}

func TestRegexp(t *testing.T) {
	tagRule := "[%d@T]%s"

	r := regexp.MustCompile(`\[(.*)@T\]`)
	// assert.Assert(t, err == nil, err)
	val := r.FindString(tagRule)
	assert.Equal(t, val, "[%d@T]")

	// r.F
	matched, err := regexp.MatchString(`\[(.*)@T\]`, tagRule)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, matched)

	loc := r.FindStringIndex(tagRule)
	assert.Equal(t, loc[0], 0)

	loc = r.FindStringIndex("25@gamer")
	assert.Equal(t, len(loc), 0)

	r = regexp.MustCompile(`\[(.*?)@T\]`)
	// r.FindStringSubmatch()
	vals := r.FindStringSubmatch("1[120@T]hostname")
	assert.Equal(t, vals[len(vals)-1], "120", val)

	r = regexp.MustCompile(`^\[(.*)@T\]`)
	vals = r.FindStringSubmatch("1[120@T]hostname")
	assert.Equal(t, len(vals), 0, vals)
}

func TestKeyTag(t *testing.T) {
	key := data.KTMem.Encode("[110@KT]25@gamer")
	assert.Equal(t, key, "[1@KT]25@gamer")

	key = data.KTDisk.Encode("")
	assert.Equal(t, key, "[2@KT]")

	key = data.KTMix.Encode("25@gamer")
	assert.Equal(t, key, "[3@KT]25@gamer")

	t.Run("equal", func(t *testing.T) {
		ekt := data.Export_KeyTag(data.KTMem)
		assert.Assert(t, ekt.Equal(data.KTMem))
		assert.Assert(t, !ekt.Equal(data.KTDisk))
		assert.Assert(t, !ekt.Equal(data.KTMix))

		ekt = data.Export_KeyTag(data.KTDisk)
		assert.Assert(t, !ekt.Equal(data.KTMem))
		assert.Assert(t, ekt.Equal(data.KTDisk))
		assert.Assert(t, !ekt.Equal(data.KTMix))

		ekt = data.Export_KeyTag(data.KTMix)
		assert.Assert(t, ekt.Equal(data.KTMem))
		assert.Assert(t, ekt.Equal(data.KTDisk))
		assert.Assert(t, ekt.Equal(data.KTMix))

		assert.Assert(t, !ekt.Equal(data.KeyTag(0)))
		assert.Assert(t, !ekt.Equal(data.KeyTag(4)))
	})

	t.Run("decode", func(t *testing.T) {
		key := data.KTMem.Encode("id")
		kt, k, err := data.Export_decode(key)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, kt, data.KTMem)
		assert.Equal(t, k, "id")

		key = "err[1@KT]id"
		kt, k, err = data.Export_decode(key)
		assert.Assert(t, err != nil)
		assert.Assert(t, kt == 0)
		assert.Equal(t, k, "")

		key = "[1@KT]"
		kt, k, err = data.Export_decode(key)
		assert.Assert(t, err != nil)
		assert.Assert(t, kt == 0)
		assert.Equal(t, k, "")

		key = "[11@KT]id"
		kt, k, err = data.Export_decode(key)
		assert.Assert(t, err != nil)
		assert.Assert(t, kt == 0)
		assert.Equal(t, k, "")

		key = "[3@KT]id"
		kt, k, err = data.Export_decode(key)
		assert.Assert(t, err == nil, nil)
		assert.Equal(t, kt, data.KTMix)
		assert.Equal(t, k, "id")
	})
}

type testKvStore struct {
	m map[string]interface{}
}

// func (ts *testKvStore)
func (ts *testKvStore) Set(k string, v interface{}) (err error) {
	if k == "" || v == nil {
		return fmt.Errorf("invalid key or value: %s : %v", k, v)
	}
	ts.m[k] = v
	return nil
}

// Get retrieves the value for the given key.
func (ts *testKvStore) Get(k string, v interface{}) (found bool, err error) {
	if k == "" || v == nil {
		return false, errors.New("invalid key or value")
	}
	stv, found := ts.m[k]
	if !found {
		return found, errors.New("cannot found value for given key")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return true, errors.New("cannot convert value to copy interface")
	}
	mrv := reflect.ValueOf(stv)
	for mrv.Kind() == reflect.Ptr {
		mrv = mrv.Elem()
	}
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	//
	rv.Set(mrv)
	return
}

// Delete deletes the stored value for the given key.
func (ts *testKvStore) Delete(k string) error {
	if k == "" {
		return errors.New("invalid key")
	}
	delete(ts.m, k)
	return nil
}

func (ts *testKvStore) Close() (err error) {
	return nil
}

func TestDataset(t *testing.T) {
	t.Run("non store in param", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil, "at least one store need to crate Dataset")
		}(t)
		data.NewDataset(nil)
	})

	onlyOneStoreTest := func(kt data.KeyTag, t *testing.T) {
		mstore := map[data.KeyTag]db.KvDB{}
		tks := &testKvStore{m: map[string]interface{}{}}
		mstore[kt] = tks
		dataset := data.NewDataset(mstore)
		err := dataset.KvDB().Set("id", 12)
		assert.Assert(t, err != nil, "invalid key")
		assert.Equal(t, len(tks.m), 0)

		err = dataset.KvDB().Set(data.KTMem.Encode("id"), 12)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(tks.m), 1)

		dataset.KvDB().Set(data.KTDisk.Encode("id"), 13)
		assert.Equal(t, len(tks.m), 1)
		assert.Equal(t, tks.m["id"], 13, "only one store")

		var val int
		found, err := dataset.KvDB().Get(data.KTMix.Encode("id"), val)
		assert.Assert(t, found)
		assert.Assert(t, err != nil)

		found, err = dataset.KvDB().Get(data.KTMix.Encode("id"), &val)
		assert.Assert(t, found)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, val, 13)

		dataset.KvDB().Set(data.KTDisk.Encode("name"), "hiank")
		assert.Equal(t, len(tks.m), 2)

		err = dataset.KvDB().Delete(data.KTDisk.Encode("id"))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(tks.m), 1)
	}

	t.Run("only memory store", func(t *testing.T) {
		onlyOneStoreTest(data.KTMem, t)
	})

	t.Run("only disk store", func(t *testing.T) {
		onlyOneStoreTest(data.KTDisk, t)
	})

	t.Run("mix store", func(t *testing.T) {
		mstore := map[data.KeyTag]db.KvDB{}
		memStore, diskStore := &testKvStore{m: map[string]interface{}{}}, &testKvStore{m: map[string]interface{}{}}
		mstore[data.KTMem], mstore[data.KTDisk] = memStore, diskStore
		dataset := data.NewDataset(mstore)
		err := dataset.KvDB().Set("id", 12)
		assert.Assert(t, err != nil, "invalid key")
		// assert.Equal(t, len(tks.m), 0)

		err = dataset.KvDB().Set(data.KTMem.Encode("id"), 12)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(diskStore.m), 0)
		assert.Equal(t, len(memStore.m), 2, "memstore will cache keytag for id")

		err = dataset.KvDB().Set(data.KTMix.Encode("id"), 13)
		assert.Assert(t, err != nil, "same key must contains same keytag")

		err = dataset.KvDB().Set(data.KTMem.Encode("id"), 14)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, memStore.m["id"], 14)
		assert.Equal(t, len(memStore.m), 2)
		assert.Equal(t, len(diskStore.m), 0)

		err = dataset.KvDB().Set(data.KTMix.Encode("name"), "h")
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(memStore.m), 4, "contains name and keytag")
		assert.Equal(t, len(diskStore.m), 1)

		var name string
		found, err := dataset.KvDB().Get(data.KTMem.Encode("name"), &name)
		assert.Assert(t, !found)
		assert.Assert(t, err != nil)

		found, err = dataset.KvDB().Get(data.KTMix.Encode("name"), &name)
		assert.Assert(t, found)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, name, "h")

		err = dataset.KvDB().Delete(data.KTMix.Encode("name"))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(memStore.m), 2)
		assert.Equal(t, len(diskStore.m), 0)

		// err = dataset.KvDB().Delete(data.)
	})
}

// func TestRedis(t *testing.T) {
// 	redis.DefaultOptions
// }
