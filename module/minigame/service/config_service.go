package service

import (
	"database/sql"
	"encoding/json"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"github.com/go-redis/redis"
	"time"
)

type ConfigService struct {
	MySql 			repository.MySqlRepository
	Cache			cache.CacheManager
	RedisService 	RedisService
	Timeout    		time.Duration
}

func NewConfigService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, timeout time.Duration) ConfigService {
	service := ConfigService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.Timeout = timeout
	return service
}

/***********************************************************
	Get prize by name
 **********************************************************/
/*
	Get Prize
*/
func (service *ConfigService) GetPrize(prizeName string) (dto.Prize, bool, error){
	// Check if it exist in Redis
	prize := dto.Prize{}
	prize, status, err := service.GetPrizeRedis(prizeName)
	if err == redis.Nil {
		err = service.RedisService.UpdateAllPrizeRedis()
		if err == nil {
			prize, status, err = service.GetPrizeRedis(prizeName)
		}
	}

	if err == nil {
		return prize, status, nil
	}

	// Update error  or get error
	prize, status, err = service.GetPrizeSQL(prizeName)
	return prize, status, err
}

/*
	Get prize redis
*/
func (service *ConfigService) GetPrizeRedis(prizeName string) (dto.Prize, bool, error){
	var prize dto.Prize
	var listPrize []dto.Prize

	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyAllPrize)
	if err == redis.Nil {
		return prize, false, err
	} else if err != nil {
		logger.Error(err.Error())
		return prize, false, err
	}
	err = json.Unmarshal([]byte(result), &listPrize)
	if err != nil {
		logger.Error(err.Error())
		return prize, false, err
	}
	for _, p := range listPrize {
		if p.Name == prizeName {
			prize = p
			break
		}
	}
	return prize, true, nil
}

/*
	Get prize SQL
*/
func (service *ConfigService) GetPrizeSQL(prizeName string) (dto.Prize, bool, error){
	prize := dto.Prize{}

	invitingUserQuery := `SELECT uuid_from_bin(Id), Value, Description FROM user_prize WHERE Name = ?;`
	invitingUserResult, err := service.MySql.DbContext.Query(invitingUserQuery, prizeName)
	if err != nil {
		service.MySql.HandleError(err)
		return prize, false, err
	}
	defer invitingUserResult.Close()

	if invitingUserResult.Next(){
		err := invitingUserResult.Scan(&prize.Id, &prize.Value, &prize.Description)
		if err != nil {
			logger.Error(err.Error())
			return prize, false, err
		}
		return prize, true, nil
	}

	return prize, false, nil
}