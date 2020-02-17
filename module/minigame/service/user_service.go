package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"github.com/go-redis/redis"
	"time"
)

type IUserService interface {
	GetWalletByUserId(ctx context.Context, userId string) (dto.User, error)
	ListTransactions(ctx context.Context, userId string, option int, pageSize int, pageIndex int) ([]dto.User, error)
	ResetRedis(ctx context.Context) error
}

type UserService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	RedisService	RedisService
	Timeout    		time.Duration
}

func NewUserService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, timeout time.Duration) IUserService {
	service := UserService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.Timeout = timeout
	return &service
}

func (service *UserService) GetWalletByUserIdSQL(ctx context.Context, userId string) (dto.User, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	var user dto.User

	userWalletQuery := `SELECT IFNULL(uuid_from_bin(game_inviting.InvitingUser), "") As InvitingUser, IFNULL(SUM(user_wallet.Value), 0) AS Wallet  
				FROM game_inviting, user_wallet 
				WHERE uuid_from_bin(user_wallet.UserId) = ? AND uuid_from_bin(game_inviting.InvitedUser) = ?;`
	userWalletResult, err := service.MySql.DbContext.Query(userWalletQuery, userId, userId)
	if err != nil {
		service.MySql.HandleError(err)
		return user, err
	}
	defer userWalletResult.Close()

	// GET read days
	now := time.Now()
	currentWeekDay := now.Weekday() - 1
	beginOfWeek := now.AddDate(0,0, - int(currentWeekDay))
	year, month, day := beginOfWeek.Date()
	lastMonday := fmt.Sprintf("%d-%d-%d", year, int(month), day)

	readDaysQuery := `SELECT dayofweek(LastUpdatedAt) 
					FROM game_read_daily 
					WHERE uuid_from_bin(UserId) = ? AND DATE(LastUpdatedAt) >= STR_TO_DATE(?,'%Y-%m-%d');`
	readDaysResult, err := service.MySql.DbContext.Query(readDaysQuery, userId, lastMonday)
	if err != nil {
		service.MySql.HandleError(err)
		return user, err
	}
	defer readDaysResult.Close()

	var wallet int
	var invitingUser string
	var isInvited bool
	var readDaysOfWeek []int

	if userWalletResult.Next(){
		err = userWalletResult.Scan(&invitingUser, &wallet)
		if err != nil {
			return user, err
		}
	}

	if invitingUser != "" {
		isInvited = true
	} else {
		isInvited = false
	}

	for readDaysResult.Next(){
		var readDayOfWeek int
		err = readDaysResult.Scan(&readDayOfWeek)
		if err != nil {
			return user, err
		}
		readDaysOfWeek = append(readDaysOfWeek, readDayOfWeek - 2)
	}

	user = dto.User{
		UserId:	userId,
		Wallet:	wallet,
		IsInvited: isInvited,
		ReadDaysOfWeek: readDaysOfWeek,
		DayOfWeek: int(currentWeekDay),
	}

	return user, nil
}

func (service *UserService) GetWalletByUserIdRedis(ctx context.Context, userId string) (dto.User, error) {
	var user dto.User
	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyUserWallet + userId)
	if err == redis.Nil {
		return user,  err
	} else if err != nil {
		logger.Error(err.Error())
		return user, err
	}

	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		logger.Error(err.Error())
		return user,  err
	}
	// GET read days
	now := time.Now()
	currentWeekDay := now.Weekday() - 1
	beginOfWeek := now.AddDate(0,0, - int(currentWeekDay))
	year, month, day := beginOfWeek.Date()
	lastMonday := fmt.Sprintf("%d-%d-%d", year, int(month), day)

	readDaysQuery := `SELECT dayofweek(LastUpdatedAt) 
					FROM game_read_daily 
					WHERE uuid_from_bin(UserId) = ? AND DATE(LastUpdatedAt) >= STR_TO_DATE(?,'%Y-%m-%d');`
	readDaysResult, err := service.MySql.DbContext.Query(readDaysQuery, userId, lastMonday)
	if err != nil {
		service.MySql.HandleError(err)
		return user, err
	}
	defer readDaysResult.Close()

	var readDaysOfWeek []int
	for readDaysResult.Next(){
		var readDayOfWeek int
		err = readDaysResult.Scan(&readDayOfWeek)
		if err != nil {
			return user, err
		}
		readDaysOfWeek = append(readDaysOfWeek, readDayOfWeek - 2)
	}

	user.ReadDaysOfWeek = readDaysOfWeek
	user.DayOfWeek 		= int(currentWeekDay)

	return user, nil
}


func (service *UserService) GetWalletByUserId(ctx context.Context, userId string) (dto.User, error) {

	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	var user dto.User
	user, err := service.GetWalletByUserIdRedis(ctx, userId)
	if err == redis.Nil {
		err = service.RedisService.UpdateUserWalletRedis(userId)
		if err == nil {
			user, err = service.GetWalletByUserIdRedis(ctx, userId)
		}
	}

	if err == nil {
		return user, nil
	}

	// Update error  or get error
	user, err = service.GetWalletByUserIdSQL(ctx, userId)

	return user, err
}


func (service *UserService) ListTransactions(ctx context.Context, userId string, option int, pageSize int, pageIndex int) ([]dto.User, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	var listTransaction []dto.User
	listTransaction, err := service.ListTransactionsRedis(ctx, userId, option, pageSize, pageIndex)
	if err == redis.Nil {
		err = service.RedisService.UpdateTransactionRedis(userId)
		if err == nil {
			listTransaction, err = service.ListTransactionsRedis(ctx, userId, option, pageSize, pageIndex)
		}
	}

	if err == nil {
		return listTransaction, nil
	}

	// Update error  or get error
	listTransaction, err = service.ListTransactionsSQL(ctx, userId, option, pageSize, pageIndex)
	return listTransaction, err
}


func (service *UserService) ListTransactionsSQL(ctx context.Context, userId string, option int, pageSize int, pageIndex int) ([]dto.User, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	var listTransactionQuery string
	switch option {
	case 0:
		listTransactionQuery = `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt FROM user_wallet, user_prize WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ? ORDER BY user_wallet.LastUpdatedAt desc LIMIT ? OFFSET ?;`
	case 1:
		listTransactionQuery = `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt FROM user_wallet, user_prize WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ? AND user_wallet.Value >= 0 ORDER BY user_wallet.LastUpdatedAt desc LIMIT ? OFFSET ?;`
	case -1:
		listTransactionQuery = `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt FROM user_wallet, user_prize WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ? AND user_wallet.Value < 0 ORDER BY user_wallet.LastUpdatedAt desc LIMIT ? OFFSET ?;`
	}

	listTransactionResult, err := service.MySql.DbContext.Query(listTransactionQuery, userId, limit, offset)
	if err != nil {
		service.MySql.HandleError(err)
		return nil, err
	}
	defer listTransactionResult.Close()

	var listUserTransaction []dto.User

	for listTransactionResult.Next(){
		var userTransaction dto.User
		err = listTransactionResult.Scan(&userTransaction.UserId, &userTransaction.Description, &userTransaction.Value, &userTransaction.LastUpdatedAt)
		if err != nil {
			return nil, err
		}
		listUserTransaction = append(listUserTransaction, userTransaction)
	}

	return listUserTransaction, nil

}

func (service *UserService) ListTransactionsRedis(ctx context.Context, userId string, option int, pageSize int, pageIndex int) ([]dto.User, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyAllTransaction + userId)
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var listTransaction []dto.User
	err = json.Unmarshal([]byte(result), &listTransaction)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var selectedTransaction []dto.User
	if option == 1 {
		for _, trans := range listTransaction {
			if trans.Value >= 0 {
				selectedTransaction = append(selectedTransaction, trans)
			}
		}
	} else if option == -1 {
		for _, trans := range listTransaction {
			if trans.Value < 0 {
				selectedTransaction = append(selectedTransaction, trans)
			}
		}
	} else {
		selectedTransaction = listTransaction
	}

	if len(selectedTransaction) >= limit + offset {
		return selectedTransaction[offset:offset+limit], nil
	} else if len(selectedTransaction) > offset &&  len(selectedTransaction) < limit + offset{
		return selectedTransaction[offset:], nil
	}else {
		return selectedTransaction[0:0], nil
	}

}


func (service *UserService) ResetRedis(ctx context.Context) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	return service.RedisService.ResetRedis()
}

