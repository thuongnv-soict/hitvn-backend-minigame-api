package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

type IPrizeService interface {
	CreatePrize(ctx context.Context, prize dto.Prize) error
	GetAllPrize(ctx context.Context, pageSize int, pageIndex int) ([]dto.Prize, error)
	UpdatePrize(ctx context.Context, prize dto.Prize) error
	DeletePrize(ctx context.Context, prizeId string) error
}

type PrizeService struct {
	MySql 			repository.MySqlRepository
	Cache			cache.CacheManager
	RedisService 	RedisService
	Timeout    		time.Duration
}

func NewPrizeService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, timeout time.Duration) IPrizeService {
	service := PrizeService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.Timeout = timeout
	return &service
}

func (service *PrizeService) CreatePrize(ctx context.Context, prize dto.Prize) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	createPrizeStatement, err := service.MySql.DbContext.Prepare(`INSERT INTO user_prize(Id, Name, Value, Description) VALUES (uuid_to_bin(?), ?, ?, ?);`)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer createPrizeStatement.Close()

	createPrizeResult, err := createPrizeStatement.Exec(util.NewUuid(), prize.Name, prize.Value, prize.Description)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	rowsAffected, err := createPrizeResult.RowsAffected()
	if err != nil {
		return errors.New("Cannot get row affected")
	}
	if rowsAffected == 0 {
		return errors.New("No row affected")
	}

	//	Update Redis
	err = service.RedisService.UpdateAllPrizeRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

/*
	Get prize SQL
*/
func (service *PrizeService) GetPrizeSQL(prizeName string) (dto.Prize, bool, error){
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

/*
	Get prize redis
*/
func (service *PrizeService) GetAllPrizeRedis(ctx context.Context, pageSize int, pageIndex int) ([]dto.Prize, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyAllPrize)
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var allPrize []dto.Prize
	err = json.Unmarshal([]byte(result), &allPrize)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if len(allPrize) >= limit + offset {
		return allPrize[offset:offset+limit], nil
	}else {
		return allPrize[offset:], nil
	}
}

/*
	Get prize sql
*/
func (service *PrizeService) GetAllPrizeSql(ctx context.Context, pageSize int, pageIndex int) ([]dto.Prize, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	getAllPrizeQuery := `SELECT uuid_from_bin(Id), Name, Value, Description, CreatedAt, LastUpdatedAt FROM user_prize ORDER BY LastUpdatedAt DESC LIMIT ? OFFSET ?`
	getAllPrizeResult, err := service.MySql.DbContext.Query(getAllPrizeQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer getAllPrizeResult.Close()

	var listPrize []dto.Prize
	for getAllPrizeResult.Next() {
		var prize dto.Prize
		err = getAllPrizeResult.Scan(&prize.Id, &prize.Name, &prize.Value, &prize.Description, &prize.CreatedAt, &prize.LastUpdatedAt)
		if err != nil {
			return nil, err
		}
		listPrize = append(listPrize, prize)
	}

	return listPrize, nil
}

/*
	Get prize
*/
func (service *PrizeService) GetAllPrize(ctx context.Context, pageSize int, pageIndex int) ([]dto.Prize, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	var listPrize []dto.Prize
	listPrize, err := service.GetAllPrizeRedis(ctx, pageSize, pageIndex)
	if err == redis.Nil {
		err = service.RedisService.UpdateAllPrizeRedis()
		if err == nil {
			listPrize, err = service.GetAllPrizeRedis(ctx, pageSize, pageIndex)
		}
	}

	if err == nil {
		return listPrize, nil
	}

	// Update error  or get error
	listPrize, err = service.GetAllPrizeSql(ctx, pageSize, pageIndex)
	return listPrize, err

}

/*
	Update Prize
 */
func (service *PrizeService) UpdatePrize(ctx context.Context, prize dto.Prize) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	updatePrizeStatement, err := service.MySql.DbContext.Prepare(`UPDATE user_prize 
																		SET Value = ?, Description = ?
																		WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	_, err = updatePrizeStatement.Exec(prize.Value, prize.Description, prize.Id)
	if err != nil {
		return err
	}
	defer updatePrizeStatement.Close()

	//	Update Redis
	err = service.RedisService.UpdateAllPrizeRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

/*
	Delete Prize
*/
func (service *PrizeService) DeletePrize(ctx context.Context, prizeId string) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	deletePrizeStatement, err := service.MySql.DbContext.Prepare(`DELETE FROM user_prize 
																		WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	_, err = deletePrizeStatement.Exec(prizeId)
	if err != nil {
		return err
	}
	defer deletePrizeStatement.Close()

	//	Update Redis
	err = service.RedisService.UpdateAllPrizeRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
