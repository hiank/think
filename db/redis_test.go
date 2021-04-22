package db_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hiank/think/db"
	"gotest.tools/v3/assert"
)

var (
	redisServerBin, _  = filepath.Abs("testdata/redis/src/redis-server")
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
	if err = exec.Command("cp", "-f", redisServerConf, dir).Run(); err != nil {
		return nil, err
	}

	baseArgs := []string{filepath.Join(dir, "redis.conf"), "--port", port, "--dir", dir}
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

func TestConnectToRedis(t *testing.T) {

	proc, err := startRedis("30211")
	if err != nil {
		t.Error(err)
		return
	}
	defer proc.Kill()

	redisConf := &db.RedisConf{}
	err = json.Unmarshal([]byte(`{
		"redis.Url": "localhost:30211",
		"redis.DB": 0,
		"redis.Password": "",
		"redis.CheckMillisecond": 300,
		"redis.TimeoutSecond": 5
	}`), redisConf)
	assert.Assert(t, err == nil, err)

	assert.Equal(t, redisConf.CheckMillisecond, 300)
	assert.Equal(t, redisConf.DB, 0)
	assert.Equal(t, redisConf.Password, "")
	assert.Equal(t, redisConf.TimeoutSecond, 5)
	assert.Equal(t, redisConf.Addr, "localhost:30211")

	ctx := context.Background()
	cli, err := db.NewVerifiedRedisCLI(ctx, redisConf)
	assert.Assert(t, err == nil, err)

	cli.Set(ctx, "testInt", 1, 0)
	val, err := cli.Get(ctx, "testInt").Result()

	assert.Assert(t, err == nil, err)
	assert.Equal(t, val, "1")
}
