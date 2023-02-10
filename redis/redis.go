package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var RedisCli *redis.Client

const RedisExpireTime = 5 * time.Second

func InitRedis(addr, password string, db int) error {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,     //"localhost:6379"
		Password: password, // no password set
		DB:       db,       // use default DB
		PoolSize: 100,      // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	_, err := RedisCli.Ping(ctx).Result()
	return err
}

func RdbExists(keys ...string) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.Exists(ctx, keys...).Result()
	return err, result
}

func RdbSet(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := RedisCli.Set(ctx, key, value, expireTime).Result()
	return err
}

func RdbGet(key string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.Get(ctx, key).Result()
	return err, result
}

func RdbHSet(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	_, err := RedisCli.HSet(ctx, key, values...).Result()
	return err
}

func RdbHGet(key, field string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.HGet(ctx, key, field).Result()
	return err, result
}

func RdbHGetAll(key string) (error, map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.HGetAll(ctx, key).Result()
	return err, result
}

func RdbHDel(key string, fields ...string) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.HDel(ctx, key, fields...).Result()
	return err, result
}

func RdbSetEx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := RedisCli.SetEX(ctx, key, value, expireTime).Result()
	return err
}

func RdbExpire(key string, expireTimeSecond uint64) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	result, err := RedisCli.Expire(ctx, key, expireTime).Result()
	return err, result
}

func RdbDel(keys ...string) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.Del(ctx, keys...).Result()
	return err, result
}

func RdbSetNx(key string, value interface{}, expireTimeSecond uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	expireTime := time.Second * time.Duration(expireTimeSecond)
	_, err := RedisCli.SetNX(ctx, key, value, expireTime).Result()
	return err
}

func RdbLPush(key string, values ...interface{}) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.LPush(ctx, key, values...).Result()
	return err, result
}

func RdbLPop(key string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.LPop(ctx, key).Result()
	return err, result
}

func RdbBLPop(waitTimeSecond uint64, keys ...string) (error, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := RedisCli.BLPop(ctx, waitTime, keys...).Result()
	return err, result
}

func RdbRPush(key string, values ...interface{}) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.RPush(ctx, key, values...).Result()
	return err, result
}

func RdbRPop(key string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.RPop(ctx, key).Result()
	return err, result
}

func RdbBRPop(waitTimeSecond uint64, keys ...string) (error, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	waitTime := time.Second * time.Duration(waitTimeSecond)
	result, err := RedisCli.BRPop(ctx, waitTime, keys...).Result()
	return err, result
}

func RdbZAdd(key string, members ...*redis.Z) (error, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.ZAdd(ctx, key, members...).Result()
	return err, result
}

func RdbZRange(key string, start, stop int64) (error, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.ZRange(ctx, key, start, stop).Result()
	return err, result
}

func RdbKeys(pattern string) (error, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.Keys(ctx, pattern).Result()
	return err, result
}

func RdbTTL(key string) (error, time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.TTL(ctx, key).Result()
	return err, result
}

func RdbDo(args ...interface{}) (error, interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), RedisExpireTime)
	defer cancel()
	result, err := RedisCli.Do(ctx, args...).Result()
	return err, result
}

func RdbDelMatchKey(cursor uint64, match string, count int64) error {
	ctx := context.Background()
	iter := RedisCli.Scan(ctx, cursor, match, count).Iterator()
	for iter.Next(ctx) {
		err := RedisCli.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
