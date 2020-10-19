package db_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/golang/glog"
	"github.com/hiank/think/db"
	"gotest.tools/assert"
)

var (
	redisServerBin, _  = filepath.Abs(filepath.Join("testdata", "redis", "src", "redis-server"))
	redisServerConf, _ = filepath.Abs(filepath.Join("testdata", "redis", "redis.conf"))
)

func testLoadGlog() {
	glog.Infoln("for repair error from --logtostderr")
}

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
	// process, err := execCmd(redisServerBin, append(baseArgs, args...)...)
	// if err != nil {
	// 	return nil, err
	// }

	// client, err := connectTo(port)
	// if err != nil {
	// 	process.Kill()
	// 	return nil, err
	// }

	// p := &redisProcess{process, client}
	// registerProcess(port, p)
	// return p, err
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
		panic(err)
	}
	defer proc.Kill()
	// db.TryRedisHub().TryMaster()
	ar, ctx := db.NewAutoRedis(context.Background(), &db.RedisConf{
		MasterURL: "localhost:30211",
		DB:        0,
		Password:  "",
	}), context.Background()
	ar.TryMaster().Set(ctx, "testInt", 1, 0)
	val, err := ar.TryMaster().Get(ctx, "testInt").Result()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, val, "1")
}
