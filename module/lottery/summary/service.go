package summary

import (
	"database/sql"
	"fmt"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"g-tech.com/module/minigame/service"
	"strings"
	"time"
)

type ISummaryLotteryService interface {
	SummaryResult(lotteryResult dto.LotteryResult) error
}

type LotterySummaryService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	ConfigService  	service.ConfigService
	RedisService	service.RedisService
	Timeout    		time.Duration
}

func NewLotterySummaryService(dbContext *sql.DB, cache cache.CacheManager, redisService service.RedisService, configService service.ConfigService, timeout time.Duration) LotterySummaryService {
	lotteryService := LotterySummaryService{}
	lotteryService.MySql.SetDbContext(dbContext)
	lotteryService.Cache = cache
	lotteryService.RedisService = redisService
	lotteryService.ConfigService = configService
	lotteryService.Timeout = timeout

	return lotteryService
}

func (service *LotterySummaryService) SummaryResult(lotteryResult dto.LotteryResult) error {
	/*
		Get Prize
	 */
	winFirstPrize, status, err := service.ConfigService.GetPrize(constant.ProgramLotteryWinFirstPrize)
	if err != nil {
		logger.Error(err.Error())
		return err
	} else {
		if status == false {
			return err
		}
	}

	/*
		Get today players
	*/
	getLotteryPlayerQuery := `SELECT uuid_from_bin(Id), uuid_from_bin(UserId), NumberSelected, Date FROM game_lottery WHERE Date = CURRENT_DATE;`
	getLotteryPlayerResult, err := service.MySql.DbContext.Query(getLotteryPlayerQuery)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer getLotteryPlayerResult.Close()

	tx, err := service.MySql.DbContext.Begin()

	var wonLotteryPlayers []dto.LotteryPlayer
	var updateWalletIdPrepare string
	for getLotteryPlayerResult.Next(){
		var lotteryPlayer dto.LotteryPlayer
		err := getLotteryPlayerResult.Scan(&lotteryPlayer.Id, &lotteryPlayer.UserId, &lotteryPlayer.NumberSelected, &lotteryPlayer.Date)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if strings.HasSuffix(lotteryResult.Special, lotteryPlayer.NumberSelected) && len(lotteryPlayer.NumberSelected) == 2{
			walletId := util.NewUuid()
			updateWalletIdPrepare += fmt.Sprintf(`UPDATE game_lottery SET WalletId = uuid_to_bin('%s') WHERE uuid_from_bin(Id) = '%s';`, walletId, lotteryPlayer.Id)
			lotteryPlayer.WalletId = walletId
			wonLotteryPlayers = append(wonLotteryPlayers, lotteryPlayer)
		}
	}

	/*
		Update user wallet for users won the lottery
	 */
	if len(wonLotteryPlayers) > 0 {
		//	Update wallet ID
		_, err = tx.Exec(updateWalletIdPrepare)
		if err != nil {
			_ = tx.Rollback()
			fmt.Println(updateWalletIdPrepare)
			logger.Error(err.Error())
			return err
		}

		// Insert User Wallet
		createWalletQuery := `INSERT INTO user_wallet(Id, UserId, PrizeId, Value) VALUES `
		for _, player := range wonLotteryPlayers{
			createWalletQuery += fmt.Sprintf(` (uuid_to_bin('%s'), uuid_to_bin('%s'), uuid_to_bin('%s'), %d),`,
				player.WalletId, player.UserId, winFirstPrize.Id, winFirstPrize.Value)
		}
		createWalletQuery = createWalletQuery[:len(createWalletQuery)-1]
		createWalletQuery += `;`
		_, err := tx.Exec(createWalletQuery)
		if err != nil {
			fmt.Println(createWalletQuery)
			logger.Error(err.Error())
			_ = tx.Rollback()
			return err
		}

		_ = tx.Commit()

		// Update Redis
		for _, player := range wonLotteryPlayers{
			err = service.RedisService.UpdateUserWalletRedis(player.UserId)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			err = service.RedisService.UpdateTransactionRedis(player.UserId)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
		}
	} else {
		_ = tx.Commit()
	}

	dt := time.Now()
	logger.Info("Date: %s\n", dt.Format("01-02-2006"))
	logger.Info("Result: %s\n", util.ToJSON(lotteryResult))
	logger.Info("Number of players won the lottery: %d\n", len(wonLotteryPlayers))

	return nil
}

