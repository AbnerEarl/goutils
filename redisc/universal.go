package redisc

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

func InitRedisUniversal(ops UniversalOptions) *UniversalClient {
	bys, _ := json.Marshal(ops)
	var ps redis.UniversalOptions
	json.Unmarshal(bys, &ps)
	rsc := redis.NewUniversalClient(&ps)
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return &UniversalClient{rsc}
}

func (rsc UniversalClient) RdbExists(keys ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Exists(ctx, keys...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbIncr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Incr(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbIncrBy(key string, value int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.IncrBy(ctx, key, value).Result()
	return result, err
}

func (rsc *UniversalClient) RdbIncrByFloat(key string, value float64) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.IncrByFloat(ctx, key, value).Result()
	return result, err
}

func (rsc *UniversalClient) RdbSet(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.Set(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *UniversalClient) RdbGet(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Get(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbHSet(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := rsc.HSet(ctx, key, values...).Result()
	return err
}

func (rsc *UniversalClient) RdbHGet(key, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HGet(ctx, key, field).Result()
	return result, err
}

func (rsc *UniversalClient) RdbHGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HGetAll(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbHDel(key string, fields ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.HDel(ctx, key, fields...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbSetEx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.SetEX(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *UniversalClient) RdbExpire(key string, expireTimeSecond uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	result, err := rsc.Expire(ctx, key, expireTime).Result()
	return result, err
}

func (rsc *UniversalClient) RdbDel(keys ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Del(ctx, keys...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbSetNx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := rsc.SetNX(ctx, key, value, expireTime).Result()
	return err
}

func (rsc *UniversalClient) RdbLPush(key string, values ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.LPush(ctx, key, values...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbLPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.LPop(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbBLPop(waitTimeSecond uint64, keys ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := rsc.BLPop(ctx, waitTime, keys...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbRPush(key string, values ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.RPush(ctx, key, values...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbRPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.RPop(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbBRPop(waitTimeSecond uint64, keys ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := rsc.BRPop(ctx, waitTime, keys...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbZAdd(key string, members []map[string]float64) (int64, error) {
	var els []*redis.Z
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

func (rsc *UniversalClient) RdbZRange(key string, start, stop int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.ZRange(ctx, key, start, stop).Result()
	return result, err
}

func (rsc *UniversalClient) RdbKeys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Keys(ctx, pattern).Result()
	return result, err
}

func (rsc *UniversalClient) RdbTTL(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.TTL(ctx, key).Result()
	return result, err
}

func (rsc *UniversalClient) RdbDo(args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	result, err := rsc.Do(ctx, args...).Result()
	return result, err
}

func (rsc *UniversalClient) RdbDelMatchKey(cursor uint64, match string, count int64) error {
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

func (rsc *UniversalClient) RdbTxPipelined(fn func(pipe Pipeliner) error) ([]Cmder, error) {
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

func (rsc *UniversalClient) RdbWatch(fn func(tx *TX) error, keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	fp := HandleFuncWatch(fn)
	err := rsc.Watch(ctx, fp, keys...)
	return err
}
