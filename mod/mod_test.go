package mod_test

import (
	"encoding/json"
	"testing"

	"github.com/hiank/think/db"
	"github.com/hiank/think/mod"
	"gotest.tools/v3/assert"
)

func TestRedis(t *testing.T) {
	t.Run("defaultValue legal", func(t *testing.T) {
		redisConf := new(db.RedisConf)
		err := json.Unmarshal([]byte(mod.Export_defaultRedisConf), redisConf)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, redisConf.CheckMillisecond, 500)
		assert.Equal(t, redisConf.DB, 0)
		assert.Equal(t, redisConf.Password, "env:REDIS_PASSWORD")
		assert.Equal(t, redisConf.TimeoutSecond, 10)

		addrMasterGetter, addrSlaveGetter := new(mod.RedisAddrMaster), new(mod.RedisAddrSlave)
		err = json.Unmarshal([]byte(mod.Export_defaultRedisAddrConf), addrMasterGetter)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, addrMasterGetter.Get().Value, "redis-master:tcp-redis")

		err = json.Unmarshal([]byte(mod.Export_defaultRedisAddrConf), addrSlaveGetter)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, addrSlaveGetter.Get().Value, "redis-slave:tcp-redis")
	})
}

func TestNats(t *testing.T) {
	t.Run("defaultValue legal", func(t *testing.T) {
		natsAddr := new(mod.NatsAddr)
		err := json.Unmarshal([]byte(mod.Export_defaultNatsConf), natsAddr)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, natsAddr.Get().Value, "nats")
	})
}
