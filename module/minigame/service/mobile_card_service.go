package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/gerror"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

type IMobileCardService interface {
	// For app
	ExchangeMobileCard(ctx context.Context, userExchange dto.UserExchange) ([]dto.MobileCard, int, error)
	GetListBoughtMobileCard(ctx context.Context, userId string, pageSize int, pageIndex int) ([]dto.MobileCard, error)

	// For web
	CreateMobileCard(ctx context.Context, mobileCard dto.MobileCard) error
	GetMobileCard(ctx context.Context, mobileCardFilter dto.MobileCardFilter, pageSize int, pageIndex int) ([]dto.MobileCard, error)
	UpdateMobileCard(ctx context.Context, mobileCard dto.MobileCard) error
	DeleteMobileCard(ctx context.Context, prizeId string) error
}

type MobileCardService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	RedisService 	RedisService
	ConfigService	ConfigService
	Timeout    		time.Duration
}

func NewMobileCardService(dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, configService ConfigService, timeout time.Duration) IMobileCardService {
	service := MobileCardService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.ConfigService = configService
	service.Timeout = timeout
	return &service
}

func (service *MobileCardService) GetListBoughtMobileCard(ctx context.Context, userId string, pageSize int, pageIndex int) ([]dto.MobileCard, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	var listBoughtMobileCard []dto.MobileCard
	listBoughtMobileCard, err := service.GetListBoughtMobileCardRedis(ctx, userId, pageSize, pageIndex)
	if err == redis.Nil {
		err = service.RedisService.UpdateBoughtMobileCardRedis(userId)
		if err == nil {
			listBoughtMobileCard, err = service.GetListBoughtMobileCardRedis(ctx, userId, pageSize, pageIndex)
		}
	}

	if err == nil {
		return listBoughtMobileCard, nil
	}

	// Update error  or get error
	listBoughtMobileCard, err = service.GetListBoughtMobileCardSQL(ctx, userId, pageSize, pageIndex)
	return listBoughtMobileCard, err
}

/*
	Get list bought mobile card Redis
*/
func (service *MobileCardService) GetListBoughtMobileCardRedis(ctx context.Context, userId string, pageSize int, pageIndex int) ([]dto.MobileCard, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyBoughtMobileCards + userId)
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var mobileCards []dto.MobileCard
	err = json.Unmarshal([]byte(result), &mobileCards)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Decode Serial and Code
	for i, _ := range mobileCards {
		mobileCards[i].Serial, err = util.DecodeMobileCard(mobileCards[i].Serial)
		if err != nil {
			return nil, err
		}

		mobileCards[i].Code, err = util.DecodeMobileCard(mobileCards[i].Code)
		if err != nil {
			return nil, err
		}
	}

	if len(mobileCards) >= limit + offset {
		return mobileCards[offset:offset+limit], nil
	} else if len(mobileCards) > offset &&  len(mobileCards) < limit + offset{
		return mobileCards[offset:], nil
	}else {
		return mobileCards[0:0], nil
	}
}

/*
	Get list bought mobile card Sql
 */
func (service *MobileCardService) GetListBoughtMobileCardSQL(ctx context.Context, userId string, pageSize int, pageIndex int) ([]dto.MobileCard, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	var mobileCardSuccessfully 	[]dto.MobileCard
	var mobileCardFailed 		[]dto.MobileCard

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	getAllBoughtMobileCardQuery := `SELECT uuid_from_bin(mobile_card.Id), mobile_card_vendor.Name, mobile_card.VendorCode, mobile_card.Serial, mobile_card.Code, mobile_card.Value, game_mobile_card.CreatedAt 
									FROM game_mobile_card, mobile_card, mobile_card_vendor 
									WHERE mobile_card_vendor.VendorCode = mobile_card.VendorCode AND game_mobile_card.MobileCardId = mobile_card.Id AND uuid_from_bin(game_mobile_card.UserId) = ?
									ORDER BY game_mobile_card.CreatedAt DESC 
									LIMIT ? OFFSET ?;`
	getAllBoughtMobileCardResult, err := service.MySql.DbContext.Query(getAllBoughtMobileCardQuery, userId, limit, offset)
	if err != nil {
		logger.Error(err.Error())
		return mobileCardFailed, err
	}
	defer getAllBoughtMobileCardResult.Close()

	for getAllBoughtMobileCardResult.Next(){
		var mobileCard dto.MobileCard
		err = getAllBoughtMobileCardResult.Scan(&mobileCard.Id, &mobileCard.Name, &mobileCard.VendorCode, &mobileCard.Serial, &mobileCard.Code, &mobileCard.Value, &mobileCard.CreatedAt)
		// Format CreatedAt
		dt,_ := time.Parse(time.RFC3339, mobileCard.CreatedAt)
		mobileCard.CreatedAt = dt.Format(constant.DateTimeLayout)

		if err != nil {
			logger.Error(err.Error())
			return mobileCardFailed, err
		}
		mobileCardSuccessfully = append(mobileCardSuccessfully, mobileCard)
	}


	// Decode Serial and Code
	for i, _ := range mobileCardSuccessfully {
		mobileCardSuccessfully[i].Serial, err = util.DecodeMobileCard(mobileCardSuccessfully[i].Serial)
		if err != nil {
			return mobileCardFailed, err
		}

		mobileCardSuccessfully[i].Code, err = util.DecodeMobileCard(mobileCardSuccessfully[i].Code)
		if err != nil {
			return mobileCardFailed, err
		}
	}

	return mobileCardSuccessfully, nil
}

/*
	Exchange mobile card
*/
func (service *MobileCardService) ExchangeMobileCard(ctx context.Context, userExchange dto.UserExchange) ([]dto.MobileCard, int, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	var mobileCardSuccessfully 	[]dto.MobileCard
	var mobileCardFailed 		[]dto.MobileCard

	// Get prize id
	prizeName := util.GetExchangeMobileCardProgramName(userExchange.Value)
	mobileCardPrize, status, err := service.ConfigService.GetPrize(prizeName)
	if err != nil {
		logger.Error(err.Error())
		return nil, 0, err
	} else {
		if status == false {
			return nil, gerror.ErrorMobileCardProgramNotFound, nil
		}
	}

	// Get user's wallet
	getUserWalletQuery := `SELECT IFNULL(SUM(Value), 0) AS Wallet  FROM user_wallet WHERE UserId = uuid_to_bin(?);`
	getUserWalletResult, err := service.MySql.DbContext.Query(getUserWalletQuery, userExchange.UserId)
	if err != nil {
		service.MySql.HandleError(err)
		return mobileCardFailed, 0, err
	}
	defer getUserWalletResult.Close()

	var wallet int
	if getUserWalletResult.Next(){
		err = getUserWalletResult.Scan(&wallet)
		if err != nil {
			return mobileCardFailed, 0, err
		}
	}


	// Check wallet enough or not?
	if wallet < - mobileCardPrize.Value * userExchange.Quantity {
		return mobileCardFailed, gerror.ErrorNotEnoughCoin, nil
	}

	getMobileCardQuery := `SELECT uuid_from_bin(mobile_card.Id) AS Id, mobile_card_vendor.Name, mobile_card.VendorCode, mobile_card.Serial, mobile_card.Code, mobile_card.Value, mobile_card.Status 
							FROM mobile_card, mobile_card_vendor
							WHERE mobile_card.VendorCode = mobile_card_vendor.VendorCode AND mobile_card_vendor.Name= ? AND mobile_card.Status = ? AND mobile_card_vendor.Status = ? AND mobile_card.Value = ?
							ORDER BY mobile_card.LastUpdatedAt ASC 
							LIMIT ?;`
	getMobileCardResult, err := service.MySql.DbContext.Query(getMobileCardQuery, userExchange.VendorName, constant.StatusMobileCardReady, constant.StatusMobileCardVendorActive, userExchange.Value, userExchange.Quantity)
	if err != nil {
		logger.Error(err.Error())
		return mobileCardFailed, 0, err
	}
	defer getMobileCardResult.Close()

	for getMobileCardResult.Next(){
		var mobileCard dto.MobileCard
		err = getMobileCardResult.Scan(&mobileCard.Id, &mobileCard.Name, &mobileCard.VendorCode, &mobileCard.Serial, &mobileCard.Code, &mobileCard.Value, &mobileCard.Status)
		if err != nil {
			logger.Error(err.Error())
			return mobileCardFailed, 0, err
		}
		mobileCardSuccessfully = append(mobileCardSuccessfully, mobileCard)
	}

	if len(mobileCardSuccessfully) == 0 {
		return mobileCardFailed, gerror.ErrorMobileCardNotExisted, err
	}

	if len(mobileCardSuccessfully) < userExchange.Quantity {
		return mobileCardFailed, gerror.ErrorNotEnoughAvailableMobileCard, err
	}

	// Decode Serial and Code
	for i, _ := range mobileCardSuccessfully {
		mobileCardSuccessfully[i].Serial, err = util.DecodeMobileCard(mobileCardSuccessfully[i].Serial)
		if err != nil {
			return mobileCardFailed, 0, err
		}

		mobileCardSuccessfully[i].Code, err = util.DecodeMobileCard(mobileCardSuccessfully[i].Code)
		if err != nil {
			return mobileCardFailed, 0, err
		}
	}


	// Add event to statistic

	// Start transaction
	tx, err := service.MySql.DbContext.Begin()

	for _, mobileCard := range mobileCardSuccessfully {

		walletId := util.NewUuid()
		// 	Add record to table user_wallet
		createMobileCardWalletStatement := `INSERT INTO user_wallet(Id, UserId, PrizeId, Value) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), ?);`
		_, err = tx.Exec(createMobileCardWalletStatement, walletId, userExchange.UserId, mobileCardPrize.Id, mobileCardPrize.Value)
		if err != nil {
			_ = tx.Rollback()
			logger.Error(err.Error())
			return mobileCardFailed, 0, err
		}

		// 	Add record to table game_mobile_card
		createMobileCardStatisticRecordStatement := `INSERT INTO game_mobile_card(Id, UserId, MobileCardId, WalletId) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`
		_, err = tx.Exec(createMobileCardStatisticRecordStatement, util.NewUuid(), userExchange.UserId, mobileCard.Id, walletId)
		if err != nil {
			_ = tx.Rollback()
			logger.Error(err.Error())
			return mobileCardFailed, 0, err
		}

		// 	Update mobile card status
		updateMobileCardStatusStatement := `UPDATE mobile_card SET status = ? WHERE Id = uuid_to_bin(?);`
		_, err = tx.Exec(updateMobileCardStatusStatement, constant.StatusMobileCardIsUsed, mobileCard.Id)
		if err != nil {
			_ = tx.Rollback()
			logger.Error(err.Error())
			return mobileCardFailed, 0, err
		}
	}

	_ = tx.Commit()

	/*
		Update Redis
	 */
	// Update Transaction
	err = service.RedisService.UpdateTransactionRedis(userExchange.UserId)
	if err != nil {
		logger.Error(err.Error())
		return mobileCardFailed, 0, err
	}

	// Update Bought Mobile Card
	err = service.RedisService.UpdateBoughtMobileCardRedis(userExchange.UserId)
	if err != nil {
		logger.Error(err.Error())
		return mobileCardFailed, 0, err
	}

	// Update user wallet
	err = service.RedisService.UpdateUserWalletRedis(userExchange.UserId)
	if err != nil {
		logger.Error(err.Error())
		return mobileCardFailed, 0, err
	}


	return mobileCardSuccessfully, 0, nil
}


/*
	Get mobile card
 */
func (service *MobileCardService) GetMobileCard(ctx context.Context, mobileCardFilter dto.MobileCardFilter, pageSize int, pageIndex int) ([]dto.MobileCard, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	limit := pageSize
	offset := (pageIndex - 1) * pageSize

	// Create query
	getMobileCardQuery := `SELECT uuid_from_bin(mobile_card.Id) AS Id, mobile_card_vendor.Name, mobile_card.VendorCode, mobile_card.Serial, mobile_card.Code, mobile_card.Value, mobile_card.Status, mobile_card.CreatedAt, mobile_card.LastUpdatedAt 
							FROM mobile_card, mobile_card_vendor
							WHERE mobile_card.VendorCode = mobile_card_vendor.VendorCode AND mobile_card_vendor.Status = 1`
	if mobileCardFilter.Name != "" {
		getMobileCardQuery += fmt.Sprintf(" AND mobile_card_vendor.Name = '%s'",  mobileCardFilter.Name)
	}
	if mobileCardFilter.Value != -1 {
		getMobileCardQuery += fmt.Sprintf(" AND mobile_card.Value = %d",  mobileCardFilter.Value)
	}
	if mobileCardFilter.Status != -1 {
		getMobileCardQuery += fmt.Sprintf(" AND mobile_card.Status = %d",  mobileCardFilter.Status)
	}
	//fmt.Println(getMobileCardQuery)
	getMobileCardQuery += " ORDER BY mobile_card.CreatedAt DESC LIMIT ? OFFSET ?;"



	getMobileCardResult, err := service.MySql.DbContext.Query(getMobileCardQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer getMobileCardResult.Close()

	var listMobileCard []dto.MobileCard
	for getMobileCardResult.Next(){
		var mobileCard dto.MobileCard
		err = getMobileCardResult.Scan(&mobileCard.Id, &mobileCard.Name, &mobileCard.VendorCode, &mobileCard.Serial, &mobileCard.Code, &mobileCard.Value, &mobileCard.Status, &mobileCard.CreatedAt, &mobileCard.LastUpdatedAt)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		listMobileCard = append(listMobileCard, mobileCard)
	}

	// Decode Serial and Code
	for i, _ := range listMobileCard {
		listMobileCard[i].Serial, err = util.DecodeMobileCard(listMobileCard[i].Serial)
		if err != nil {
			return nil, err
		}

		listMobileCard[i].Code, err = util.DecodeMobileCard(listMobileCard[i].Code)
		if err != nil {
			return nil, err
		}
	}

	return listMobileCard, nil
}

/*
	Create mobile card
*/
func (service *MobileCardService) CreateMobileCard(ctx context.Context, mobileCard dto.MobileCard) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	createMobileCardStatement, err := service.MySql.DbContext.Prepare(`INSERT INTO mobile_card(Id, VendorCode, Serial, Code, Value, Status) VALUES (uuid_to_bin(?), ?, ?, ?, ?, ?);`)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer createMobileCardStatement.Close()

	mobileCard.Serial, _ 	= util.EncodeMobileCard(mobileCard.Serial)
	mobileCard.Code, _ 		= util.EncodeMobileCard(mobileCard.Code)

	createMobileCardResult, err := createMobileCardStatement.Exec(util.NewUuid(), mobileCard.VendorCode, mobileCard.Serial, mobileCard.Code, mobileCard.Value, mobileCard.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	rowsAffected, err := createMobileCardResult.RowsAffected()
	if err != nil {
		return errors.New("Cannot get row affected")
	}
	if rowsAffected == 0 {
		return errors.New("No row affected")
	}

	return nil
}

/*
	Update Mobile Card
*/
func (service *MobileCardService) UpdateMobileCard(ctx context.Context, mobileCard dto.MobileCard) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	mobileCard.Serial, _ 	= util.EncodeMobileCard(mobileCard.Serial)
	mobileCard.Code, _ 		= util.EncodeMobileCard(mobileCard.Code)
	updateMobileCardStatement, err := service.MySql.DbContext.Prepare(`UPDATE mobile_card 
																			SET VendorCode = ?, Serial = ?,
																				Code = ?, Value = ?, Status = ?
																			WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	_, err = updateMobileCardStatement.Exec(mobileCard.VendorCode, mobileCard.Serial, mobileCard.Code, mobileCard.Value, mobileCard.Status, mobileCard.Id)
	if err != nil {
		return err
	}
	defer updateMobileCardStatement.Close()

	return nil
}

/*
	Delete Mobile Card
*/
func (service *MobileCardService) DeleteMobileCard(ctx context.Context, mobileCardId string) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	deleteMobileCardStatement, err := service.MySql.DbContext.Prepare(`DELETE FROM mobile_card 
																		WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	_, err = deleteMobileCardStatement.Exec(mobileCardId)
	if err != nil {
		return err
	}
	defer deleteMobileCardStatement.Close()

	return nil
}

/*
	Update Redis
*/
//func (summary *MobileCardService) UpdateRedis(userId string) error{
//	// 1. Update list transaction
//
//	listTransactionQuery 	:= `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt
//								FROM user_wallet, user_prize
//								WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ?
//								ORDER BY user_wallet.LastUpdatedAt DESC;`
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
//
//	// 2. Update bought mobile cards
//	var mobileCardSuccessfully 	[]dto.MobileCard
//
//	getAllBoughtMobileCardQuery := `SELECT uuid_from_bin(mobile_card.Id), mobile_card_vendor.Name, mobile_card.VendorCode, mobile_card.Serial, mobile_card.Code, mobile_card.Value, game_mobile_card.CreatedAt
//									FROM game_mobile_card, mobile_card, mobile_card_vendor
//									WHERE mobile_card_vendor.VendorCode = mobile_card.VendorCode AND game_mobile_card.MobileCardId = mobile_card.Id AND uuid_from_bin(game_mobile_card.UserId) = ?
//									ORDER BY game_mobile_card.CreatedAt DESC;`
//	getAllBoughtMobileCardResult, err := summary.MySql.DbContext.Query(getAllBoughtMobileCardQuery, userId)
//	if err != nil {
//		logger.Error(err.Error())
//		return err
//	}
//	defer getAllBoughtMobileCardResult.Close()
//
//	for getAllBoughtMobileCardResult.Next(){
//		var mobileCard dto.MobileCard
//		err = getAllBoughtMobileCardResult.Scan(&mobileCard.Id, &mobileCard.Name, &mobileCard.VendorCode, &mobileCard.Serial, &mobileCard.Code, &mobileCard.Value, &mobileCard.CreatedAt)
//		// Format CreatedAt
//		dt,_ := time.Parse(time.RFC3339, mobileCard.CreatedAt)
//		mobileCard.CreatedAt = dt.Format(constant.DateTimeLayout)
//
//		if err != nil {
//			logger.Error(err.Error())
//			return err
//		}
//		mobileCardSuccessfully = append(mobileCardSuccessfully, mobileCard)
//	}
//
//	err = summary.RedisClient.Set(constant.RedisPrefixKeyBoughtMobileCards + userId, util.ToJSON(mobileCardSuccessfully), 0).Err()
//	if err != nil {
//		logger.Error(err.Error())
//		return err
//	}
//
//	return nil
//}
