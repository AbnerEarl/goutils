package redisc

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func InitRedisCluster(addrs []string, db, poolSize, idleConns int, username, password string) *RedisClusterCli {
	if db < 0 {
		db = 0
	}
	if poolSize < 1 {
		poolSize = 100
	}
	rsc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,     // []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}
		Username:     username,  // no username set
		Password:     password,  // no password set
		PoolSize:     poolSize,  // 连接池大小
		MinIdleConns: idleConns, // 最小连接大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return &RedisClusterCli{rsc}
}

func InitRedisFailoverCluster(masterName string, sentinelAddrs []string, isRandom bool, db, poolSize, idleConns int, username, password, sentinelUsername, sentinelPassword string) *RedisClusterCli {
	if db < 0 {
		db = 0
	}
	if poolSize < 1 {
		poolSize = 100
	}
	ops := redis.FailoverOptions{
		MasterName:       masterName,       // "master"
		SentinelAddrs:    sentinelAddrs,    // []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"
		Username:         username,         // no username set
		Password:         password,         // no password set
		SentinelUsername: sentinelUsername, // no username set
		SentinelPassword: sentinelPassword, // no password set
		DB:               db,               // use default DB
		PoolSize:         poolSize,         // 连接池大小
		MinIdleConns:     idleConns,        // 最小连接大小
	}
	if isRandom {
		ops.RouteRandomly = true
	} else {
		ops.RouteByLatency = true
	}
	rsc := redis.NewFailoverClusterClient(&ops)

	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return &RedisClusterCli{rsc}
}

func (rsc *RedisClusterCli) RdbExists(keys ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Exists(ctx, keys...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbSet(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.Set(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *RedisClusterCli) RdbGet(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Get(ctx, key).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbHSet(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.HSet(ctx, key, values...).Result()
	return err
}

func (rsc *RedisClusterCli) RdbHGet(key, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HGet(ctx, key, field).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbHGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HGetAll(ctx, key).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbHDel(key string, fields ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HDel(ctx, key, fields...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbSetEx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.SetEX(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *RedisClusterCli) RdbExpire(key string, expireTimeSecond uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	result, err := rsc.Expire(ctx, key, expireTime).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbDel(keys ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Del(ctx, keys...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbSetNx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.SetNX(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *RedisClusterCli) RdbLPush(key string, values ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.LPush(ctx, key, values...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbLPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.LPop(ctx, key).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbBLPop(waitTimeSecond uint64, keys ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := rsc.BLPop(ctx, waitTime, keys...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbRPush(key string, values ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.RPush(ctx, key, values...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbRPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.RPop(ctx, key).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbBRPop(waitTimeSecond uint64, keys ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := rsc.BRPop(ctx, waitTime, keys...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbZAdd(key string, members []map[string]float64) (int64, error) {
	els := []*redis.Z{}
	for _, m := range members {
		for k, v := range m {
			els = append(els, &redis.Z{
				Score:  v,
				Member: k,
			})
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.ZAdd(ctx, key, els...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbZRange(key string, start, stop int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.ZRange(ctx, key, start, stop).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbKeys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Keys(ctx, pattern).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbTTL(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.TTL(ctx, key).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbDo(args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Do(ctx, args...).Result()
	return result, err
}

func (rsc *RedisClusterCli) RdbDelMatchKey(cursor uint64, match string, count int64) error {
	ctx := context.Background()
	iter := rsc.Scan(ctx, cursor, match, count).Iterator()
	for iter.Next(ctx) {
		err := rsc.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (rsc *RedisClusterCli) RdbTxPipelined(fn func(pipe Pipeliner) error) ([]Cmder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	fp := HandleFuncPipe(fn)
	cmders, err := rsc.TxPipelined(ctx, fp)
	var result []Cmder
	for _, cmd := range cmders {
		result = append(result, Cmder(cmd))
	}
	return result, err
}

func (rsc *RedisClusterCli) RdbWatch(fn func(tx *TX) error, keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	fp := HandleFuncWatch(fn)
	err := rsc.Watch(ctx, fp, keys...)
	return err
}
