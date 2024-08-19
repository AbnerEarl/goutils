/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/8/16 09:26
 * @desc: about the role of class.
 */

package web

import (
	"crypto/rand"
	"github.com/AbnerEarl/goutils/redisc"
	"math/big"
)

const (
	MAX_RANDOM_NUMBER             = 256
	MIN_RANDOM_NUMBER             = 0
	SALT_ONE_LENGTH               = 8
	SALT_TWO_LENGTH               = 8
	INCREASE_COUNTER_LENGTH       = 47
	HIGHEST_NUMBER                = "0"
	BINARY_PREFIX                 = "0b"
	UUID_NUMBER_LENGTH            = 63
	SALT_ONE_SHIFT                = SALT_TWO_LENGTH
	INCREASE_COUNTER_SHIFT        = SALT_TWO_LENGTH + SALT_ONE_LENGTH
	AUTO_ID_GENERATOR_COUNTER_KEY = "auto_id_generator_counter_key"
)

func GenAutoId(client *redisc.RedisCli) uint64 {
	saltOne, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
	saltTwo, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
	incrValue, _ := client.RdbIncr(AUTO_ID_GENERATOR_COUNTER_KEY)
	result := (incrValue << INCREASE_COUNTER_SHIFT) | (saltOne.Int64() << SALT_ONE_SHIFT) | saltTwo.Int64()
	return uint64(result)
}

func GenAutoIds(n uint64, client *redisc.RedisCli) []uint64 {
	result := []uint64{}

	for i := uint64(0); i < n; i++ {
		saltOne, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
		saltTwo, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
		incrValue, _ := client.RdbIncr(AUTO_ID_GENERATOR_COUNTER_KEY)
		id := (incrValue << INCREASE_COUNTER_SHIFT) | (saltOne.Int64() << SALT_ONE_SHIFT) | saltTwo.Int64()
		result = append(result, uint64(id))
	}
	return result
}

func GenAutoIdByClu(client *redisc.RedisClusterCli) uint64 {
	saltOne, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
	saltTwo, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
	incrValue, _ := client.RdbIncr(AUTO_ID_GENERATOR_COUNTER_KEY)
	result := (incrValue << INCREASE_COUNTER_SHIFT) | (saltOne.Int64() << SALT_ONE_SHIFT) | saltTwo.Int64()
	return uint64(result)
}

func GenAutoIdsByClu(n uint64, client *redisc.RedisClusterCli) []uint64 {
	result := []uint64{}

	for i := uint64(0); i < n; i++ {
		saltOne, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
		saltTwo, _ := rand.Int(rand.Reader, big.NewInt(MAX_RANDOM_NUMBER))
		incrValue, _ := client.RdbIncr(AUTO_ID_GENERATOR_COUNTER_KEY)
		id := (incrValue << INCREASE_COUNTER_SHIFT) | (saltOne.Int64() << SALT_ONE_SHIFT) | saltTwo.Int64()
		result = append(result, uint64(id))
	}
	return result
}
