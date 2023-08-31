package redisc

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func InitRedisSentinel(addr string, db, poolSize, idleConns int, username, password string) *SentinelClient {
	if db < 0 {
		db = 0
	}
	if poolSize < 1 {
		poolSize = 100
	}
	rsc := redis.NewSentinelClient(&redis.Options{
		Addr:         addr,      // "localhost:6379"
		Username:     username,  // no username set
		Password:     password,  // no password set
		DB:           db,        // use default DB
		PoolSize:     poolSize,  // 连接池大小
		MinIdleConns: idleConns, // 最小连接大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return &SentinelClient{rsc}
}

func (rsc *SentinelClient) RdbSet(name, option, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.Set(ctx, name, option, value).Result()
	return err
}

func (rsc *SentinelClient) RdbGetMasterAddrByName(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.GetMasterAddrByName(ctx, key).Result()
	return result, err
}

func (rsc *SentinelClient) RdbFailover(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Failover(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbCkQuorum(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.CkQuorum(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbFlushConfig() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.FlushConfig(ctx).Result()
	return result, err
}

func (rsc *SentinelClient) RdbMaster(name string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Master(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbMasters(key string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Masters(ctx).Result()
	return result, err
}

func (rsc *SentinelClient) RdbMonitor(name, ip, port, quorum string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Monitor(ctx, name, ip, port, quorum).Result()
	return result, err
}

func (rsc *SentinelClient) RdbProcess(cmder Cmder) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := rsc.Process(ctx, cmder)
	return err
}

func (rsc *SentinelClient) RdbPSubscribe(channels ...string) *redis.PubSub {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result := rsc.PSubscribe(ctx, channels...)
	return result
}

func (rsc *SentinelClient) RdbRemove(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Remove(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbSentinels(name string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Sentinels(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbSlaves(name string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Slaves(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbReset(name string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Reset(ctx, name).Result()
	return result, err
}

func (rsc *SentinelClient) RdbSubscribe(name string) *redis.PubSub {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result := rsc.Subscribe(ctx, name)
	return result
}
