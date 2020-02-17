package service

import (
	"context"
	"database/sql"
	"errors"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/gerror"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"time"
)

type ILotteryService interface {
	CreateLotteryNumber(ctx context.Context, lotteryNumber dto.LotteryPlayer) (int, error)
	GetSelectedNumbers(ctx context.Context, userId string) ([]dto.LotteryPlayer, error)
}

type LotteryService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	RedisService	RedisService
	ConfigService 	ConfigService
	Timeout    		time.Duration
}

func NewLotteryService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, configService ConfigService, timeout time.Duration) ILotteryService {
	service := LotteryService{}
	service.Cache = cache
	service.RedisService = redisService
	service.ConfigService = configService
	service.MySql.SetDbContext(dbContext)
	service.Timeout = timeout
	return &service
}

func (service *LotteryService) CreateLotteryNumber(ctx context.Context, lotteryPlayer dto.LotteryPlayer) (int, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()
	/*
		Check time add number
	 */
	if time.Now().Hour() >= 18 {
		return gerror.ErrorLotteryTimeUp, nil
	}

	/*
		Check selected number
	 */
	getNumberOfSelectedQuery := `SELECT NumberSelected
									FROM game_lottery
									WHERE uuid_from_bin(UserId) = ? AND Date = CURRENT_DATE ;`
	getNumberOfSelectedResult, err := service.MySql.DbContext.Query(getNumberOfSelectedQuery, lotteryPlayer.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}
	defer getNumberOfSelectedResult.Close()

	var selectedNumbers []string
	for getNumberOfSelectedResult.Next(){
		var selectedNumber string
		err = getNumberOfSelectedResult.Scan(&selectedNumber)
		if err != nil {
			logger.Error(err.Error())
			return 0, err
		}
		selectedNumbers = append(selectedNumbers, selectedNumber)
	}

	if len(selectedNumbers) >= constant.DefaultMaximumSelectedLotteryNumbers {
		return gerror.ErrorLotteryExceedNumberOfSelected, nil
	}

	for _, selectedNumber := range selectedNumbers {
		if selectedNumber == lotteryPlayer.NumberSelected {
			return gerror.ErrorLotteryDuplicatedSelectedNumber, nil
		}
	}

	/*
		Create Selected Lottery Numbers
	 */
	createLotteryNumberStatement, err := service.MySql.DbContext.Prepare(`INSERT INTO game_lottery(Id, UserId, NumberSelected, Date) VALUES (uuid_to_bin(?), uuid_to_bin(?), ?, CURRENT_DATE);`)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}
	defer createLotteryNumberStatement.Close()

	createMobileCardResult, err := createLotteryNumberStatement.Exec(util.NewUuid(), lotteryPlayer.UserId, lotteryPlayer.NumberSelected)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	rowsAffected, err := createMobileCardResult.RowsAffected()
	if err != nil {
		return 0, errors.New("Cannot get row affected")
	}
	if rowsAffected == 0 {
		return 0, errors.New("No row affected")
	}

	return 0, nil
}


/*
	GET selected numbers
*/
func (service *LotteryService) GetSelectedNumbers(ctx context.Context, userId string) ([]dto.LotteryPlayer, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	getSelectedNumbersQuery := `SELECT uuid_from_bin(Id), uuid_from_bin(UserId), NumberSelected, Date
									FROM game_lottery
									WHERE uuid_from_bin(UserId) = ? AND Date = CURRENT_DATE 
									ORDER BY CreatedAt ASC;`
	getSelectedNumbersResult, err := service.MySql.DbContext.Query(getSelectedNumbersQuery, userId)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer getSelectedNumbersResult.Close()

	var lotteryPlayers []dto.LotteryPlayer
	for getSelectedNumbersResult.Next() {
		var lotteryPlayer dto.LotteryPlayer
		err = getSelectedNumbersResult.Scan(&lotteryPlayer.Id, &lotteryPlayer.UserId, &lotteryPlayer.NumberSelected, &lotteryPlayer.Date)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		lotteryPlayers = append(lotteryPlayers, lotteryPlayer)
	}

	return lotteryPlayers, nil
}