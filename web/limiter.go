/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/7/26 14:35
 * @desc: about the role of class.
 */

package web

import (
	"encoding/json"
	"fmt"
	"github.com/AbnerEarl/goutils/redisc"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

var funnels = sync.Map{}

type funnelRateLimiter struct {
	Capacity    int64
	LeakingRate float64
	LeftQuota   int64
	LeakingTs   int64
}

func (frl *funnelRateLimiter) makeSpace() {
	nowTs := time.Millisecond.Milliseconds()
	deltaTs := nowTs - frl.LeakingTs
	deltaQuota := int64(float64(deltaTs) * frl.LeakingRate)
	if deltaQuota < 0 {
		frl.LeftQuota = frl.Capacity
		frl.LeakingTs = nowTs
		return
	} else if deltaQuota < 1 {
		return
	}
	frl.LeftQuota += deltaQuota
	frl.LeakingTs = nowTs
	if frl.LeftQuota > frl.Capacity {
		frl.LeftQuota = frl.Capacity
	}
}

func (frl *funnelRateLimiter) watering(quota int64) bool {
	frl.makeSpace()
	if frl.LeftQuota > quota {
		frl.LeftQuota -= quota
		return true
	}
	return false
}

func IsAllowedByMap(userId, actionKey string, capacity int64, leakingRate float64) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	funnel, ok := funnels.Load(key)
	if !ok {
		frl := funnelRateLimiter{
			Capacity:    capacity,
			LeakingRate: leakingRate,
			LeftQuota:   capacity,
			LeakingTs:   time.Millisecond.Milliseconds(),
		}
		funnels.Store(key, frl)
		return frl.watering(1)
	}
	frl := funnel.(funnelRateLimiter)
	return frl.watering(1)
}

func IsAllowedByRedis(userId, actionKey string, capacity int64, leakingRate float64, client *redisc.RedisCli) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	result, err := client.RdbGet(key)
	expTime := uint64(float64(capacity)*leakingRate) + 1
	if err != nil {
		frl := funnelRateLimiter{
			Capacity:    capacity,
			LeakingRate: leakingRate,
			LeftQuota:   capacity,
			LeakingTs:   time.Millisecond.Milliseconds(),
		}
		client.RdbSet(key, frl, expTime)
		return frl.watering(1)
	}
	client.RdbExpire(key, expTime)
	frl := funnelRateLimiter{}
	json.Unmarshal([]byte(result), &frl)
	return frl.watering(1)
}

func IsAllowedByRedisCluster(userId, actionKey string, capacity int64, leakingRate float64, client *redisc.RedisClusterCli) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	result, err := client.RdbGet(key)
	expTime := uint64(float64(capacity)*leakingRate) + 1
	if err != nil {
		frl := funnelRateLimiter{
			Capacity:    capacity,
			LeakingRate: leakingRate,
			LeftQuota:   capacity,
			LeakingTs:   time.Millisecond.Milliseconds(),
		}
		client.RdbSet(key, frl, expTime)
		return frl.watering(1)
	}
	client.RdbExpire(key, expTime)
	frl := funnelRateLimiter{}
	json.Unmarshal([]byte(result), &frl)
	return frl.watering(1)
}

func IsLimitedByMap(userId, actionKey string, capacity int, leakingRate float64) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	funnel, ok := funnels.Load(key)
	if !ok {
		limiter := rate.NewLimiter(rate.Limit(leakingRate), capacity)
		funnels.Store(key, limiter)
		return limiter.Allow()
	}
	limiter := funnel.(rate.Limiter)
	return limiter.Allow()
}

func IsLimitedByRedis(userId, actionKey string, capacity int, leakingRate float64, client *redisc.RedisCli) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	result, err := client.RdbGet(key)
	expTime := uint64(float64(capacity)*leakingRate) + 1
	if err != nil {
		limiter := rate.NewLimiter(rate.Limit(leakingRate), capacity)
		client.RdbSet(key, limiter, expTime)
		return limiter.Allow()
	}
	client.RdbExpire(key, expTime)
	limiter := rate.Limiter{}
	json.Unmarshal([]byte(result), &limiter)
	return limiter.Allow()
}

func IsLimitedByRedisCluster(userId, actionKey string, capacity int, leakingRate float64, client *redisc.RedisClusterCli) bool {
	key := fmt.Sprintf("%s:%s", userId, actionKey)
	result, err := client.RdbGet(key)
	expTime := uint64(float64(capacity)*leakingRate) + 1
	if err != nil {
		limiter := rate.NewLimiter(rate.Limit(leakingRate), capacity)
		client.RdbSet(key, limiter, expTime)
		return limiter.Allow()
	}
	client.RdbExpire(key, expTime)
	limiter := rate.Limiter{}
	json.Unmarshal([]byte(result), &limiter)
	return limiter.Allow()
}
