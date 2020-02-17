package service

import (
	"context"
	"database/sql"
	"fmt"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/gerror"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"time"
)

type IReadDailyService interface {
	CreateNewReadDaily(ctx context.Context, user dto.User) (int, error)
}

type ReadDailyService struct {
	MySql 			repository.MySqlRepository
	RedisService 	RedisService
	Cache 			cache.CacheManager
	Timeout    		time.Duration
}

func NewReadDailyService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, timeout time.Duration) IReadDailyService {
	service := ReadDailyService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.Timeout = timeout
	return &service
}

func (service *ReadDailyService) CreateNewReadDaily (ctx context.Context, user dto.User) (int, error) {
	/*
		Check user has received coin today or not?
	*/
	checkReceivedQuery := `SELECT Id FROM game_read_daily WHERE uuid_from_bin(UserId) = ? AND DATE(LastUpdatedAt) = CURRENT_DATE;`
	checkReceivedResult, err := service.MySql.DbContext.Query(checkReceivedQuery, user.UserId)
	if err != nil {
		service.MySql.HandleError(err)
		return 0, err
	}
	defer checkReceivedResult.Close()

	if checkReceivedResult.Next(){
		return gerror.ErrorReadDailyUserHasReceivedCoinToday, nil
	}

	/*
		Get Prize
	*/
	prizeId := ""
	prizeValue := 0

	invitingUserQuery := `SELECT uuid_from_bin(Id), Value FROM user_prize WHERE Name = ?;`
	invitingUserResult, err := service.MySql.DbContext.Query(invitingUserQuery, constant.ProgramReadHITDaily)
	if err != nil {
		service.MySql.HandleError(err)
		return 0, err
	}
	defer invitingUserResult.Close()

	if invitingUserResult.Next(){
		err := invitingUserResult.Scan(&prizeId, &prizeValue)
		if err != nil {
			logger.Error(err.Error())
			return 0, err
		}
	} else {
		return gerror.ErrorReadDailyProgramNotFound, nil
	}

	// Start transaction
	tx, err := service.MySql.DbContext.Begin()

	walletId := util.NewUuid()
	createReadDailyStatement := `INSERT INTO user_wallet(Id, UserId, PrizeId, Value) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), ?);`

	_, err = tx.Exec(createReadDailyStatement, walletId, user.UserId, prizeId, user.Value)
	if err != nil {
		fmt.Println(err.Error())
		_ = tx.Rollback()

		return 0, err
	}

	createWalletStatement := `INSERT INTO game_read_daily(Id, UserId, WalletId) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`
	_, err = tx.Exec(createWalletStatement, util.NewUuid(), user.UserId, walletId)
	if err != nil {
		fmt.Println(err.Error())
		_ = tx.Rollback()
		return 0, err
	}

	_ = tx.Commit()

	/*
		Update Redis
	*/
	// Update transaction
	err = service.RedisService.UpdateTransactionRedis(user.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	// Update user wallet
	err = service.RedisService.UpdateUserWalletRedis(user.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	return 0, nil
}


func (service *ReadDailyService) CreateNewReadDailyV2(ctx context.Context, user dto.User) (int, error) {
	/*
		Check user has received coin today or not?
	*/
	checkReceivedQuery := `SELECT Id FROM game_read_daily WHERE uuid_from_bin(UserId) = ? AND DATE(LastUpdatedAt) = CURRENT_DATE;`
	checkReceivedResult, err := service.MySql.DbContext.Query(checkReceivedQuery, user.UserId)
	if err != nil {
		service.MySql.HandleError(err)
		return 0, err
	}
	defer checkReceivedResult.Close()

	if checkReceivedResult.Next(){
		return gerror.ErrorReadDailyUserHasReceivedCoinToday, nil
	}

	/*
		Get Prize
	*/
	prizeId := ""
	prizeValue := 0

	invitingUserQuery := `SELECT uuid_from_bin(Id), Value FROM user_prize WHERE Name = ?;`
	invitingUserResult, err := service.MySql.DbContext.Query(invitingUserQuery, constant.ProgramReadHITDaily)
	if err != nil {
		service.MySql.HandleError(err)
		return 0, err
	}
	defer invitingUserResult.Close()

	if invitingUserResult.Next(){
		err := invitingUserResult.Scan(&prizeId, &prizeValue)
		if err != nil {
			logger.Error(err.Error())
			return 0, err
		}
	} else {
		return gerror.ErrorReadDailyProgramNotFound, nil
	}

	// Start transaction
	tx, err := service.MySql.DbContext.Begin()

	walletId := util.NewUuid()
	createReadDailyStatement := `INSERT INTO user_wallet(Id, UserId, PrizeId, Value) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), ?);`

	_, err = tx.Exec(createReadDailyStatement, walletId, user.UserId, prizeId, user.Value)
	if err != nil {
		fmt.Println(err.Error())
		_ = tx.Rollback()

		return 0, err
	}

	createWalletStatement := `INSERT INTO game_read_daily(Id, UserId, WalletId) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`
	_, err = tx.Exec(createWalletStatement, util.NewUuid(), user.UserId, walletId)
	if err != nil {
		fmt.Println(err.Error())
		_ = tx.Rollback()
		return 0, err
	}

	_ = tx.Commit()

	/*
		Update Redis
	*/
	// Update transaction
	err = service.RedisService.UpdateTransactionRedis(user.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	// Update user wallet
	err = service.RedisService.UpdateUserWalletRedis(user.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	return 0, nil
}
//
//func (summary *ReadDailyService) UpdateRedis(userId string) error{
//
//	listTransactionQuery 	:= `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt FROM user_wallet, user_prize WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ? ORDER BY user_wallet.LastUpdatedAt DESC;`
//
//	listTransactionResult, err := summary.MySql.DbContext.Query(listTransactionQuery, userId)
//	if err != nil {
//		summary.MySql.HandleError(err)
//		return err
//	}
//	defer listTransactionResult.Close()
//
//	var listUserTransaction []dto.User
//
//	for listTransactionResult.Next(){
//		var userTransaction dto.User
//		err = listTransactionResult.Scan(&userTransaction.UserId, &userTransaction.Description, &userTransaction.Value, &userTransaction.LastUpdatedAt)
//		if err != nil {
//			return err
//		}
//		// Format CreatedAt
//		dt,_ := time.Parse(time.RFC3339, userTransaction.LastUpdatedAt)
//		userTransaction.LastUpdatedAt = dt.Format(constant.DateTimeLayout)
//
//		listUserTransaction = append(listUserTransaction, userTransaction)
//	}
//
//	err = summary.RedisClient.Set(constant.RedisPrefixKeyAllTransaction + userId, util.ToJSON(listUserTransaction), 0).Err()
//	if err != nil {
//		logger.Error(err.Error())
//		return err
//	}
//
//	return nil
//}
