package cache

import (
	"encoding/json"
	"fmt"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/util"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type CacheManager struct {
	Client *redis.Client
}

/**
 * Initializes cache
 */
func (manager *CacheManager) Init(host string, poolSize int, minIdleConns int) {

	manager.Client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: poolSize,
		MinIdleConns: minIdleConns,
	})

}
/**
 * Delete item
 */
func (manager *CacheManager) DeleteItem(key string) int64 {
	err := manager.Client.Del(key).Val()
	return err
}
/*
* Push item to queue list
*/
func (manager *CacheManager) PushItem(queueKey string, value interface{}) error{
	err := manager.Client.RPush(queueKey, value).Err()
	return err
}
/*
* Pop item from queue list
*/
func (manager *CacheManager) PopItem(queueKey string) (interface{}, error){
	return manager.Client.LPop(queueKey).Result()
}
/*
* Get Score of Member in Sorted List
*/
func (manager *CacheManager) ZGetScore(key string, member string) float64{
	return manager.Client.ZScore(key, member).Val()
}
/*
*
*/
func (manager *CacheManager)ZRevRangeByScore(key string, startScore float64, intPageSize int64) ([]redis.Z, error){

	var max string
	if startScore != -1{
	max = fmt.Sprintf("(%.f", startScore)
	}else{
		max = "+inf"
	}
	return manager.Client.ZRevRangeByScoreWithScores(key, redis.ZRangeBy{
		Max: max,
		Offset: 0,
		Count: intPageSize,
	}).Result()
}
/*
* Get Count of record by member in Sorted List
*/
func (manager *CacheManager)ZGetCountByMember(key string, member string) (int64){
	var score = manager.Client.ZScore(key, member).Val()

	return manager.Client.ZCount(key, util.FloatToString(score), "+inf").Val()
}
/*
* Get a item from cache
*/
func (manager *CacheManager) Get(key string) string {
	value, err := manager.Client.Get(key).Result()
	if err != nil && !strings.Contains(err.Error(), "redis: nil"){
		logger.Error(err.Error())
	}

	return value
}
/*
* Get a item from cache with error
 */
func (manager *CacheManager) GetWithError(key string) (string, error){
	value, err := manager.Client.Get(key).Result()
	if err != nil && !strings.Contains(err.Error(), "redis: nil"){
		logger.Error(err.Error())
	}

	return value, err
}

/*
* Set a item to cache
*/
func (manager *CacheManager) Set(key string, object interface{}, expireIn time.Duration){
	out, err := json.Marshal(object)
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		logger.Error(err.Error())
	}

	err = manager.Client.Set(key, out, expireIn).Err()
	if err != nil {
		logger.Error(err.Error())
	}
}

/*
* Set a item to cache with error
 */
func (manager *CacheManager) SetWithError(key string, object interface{}, expireIn time.Duration) error{
	out, err := json.Marshal(object)
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		logger.Error(err.Error())
		return err
	}

	err = manager.Client.Set(key, out, expireIn).Err()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (manager *CacheManager) RPush(key, value string) {
	err := manager.Client.RPush(key, value).Err()
	if err != nil && !strings.Contains(err.Error(), "redis: nil"){
		logger.Error(err.Error())
	}
}

func (manager *CacheManager) LPop(key string) interface{}{
	return manager.Client.LPop(key)
}

func (manager *CacheManager) LGetFirst(key string) interface{}{
	return manager.Client.LRange(key, 0, 0)
}

func (manager *CacheManager) LGetAll(key string) []string{
	return manager.Client.LRange(key, 0, -1).Val()
}