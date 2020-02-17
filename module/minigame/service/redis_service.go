package service

import (
	"database/sql"
	"fmt"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"time"
)

type RedisService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	Timeout    		time.Duration
}

func NewRedisService (dbContext *sql.DB, cache cache.CacheManager, timeout time.Duration) RedisService {
	service := RedisService{}
	service.Cache = cache
	service.MySql.SetDbContext(dbContext)
	service.Timeout = timeout
	return service
}

/*
	Update transaction
*/
func (service *RedisService) UpdateTransactionRedis(userId string) error{

	listTransactionQuery 	:= `SELECT uuid_from_bin(user_wallet.UserId) AS UserId, user_prize.Description, user_wallet.Value, user_wallet.LastUpdatedAt 
								FROM user_wallet, user_prize 
								WHERE user_wallet.PrizeId = user_prize.Id AND uuid_from_bin(user_wallet.UserId) = ? 
								ORDER BY user_wallet.LastUpdatedAt DESC;`
	listTransactionResult, err := service.MySql.DbContext.Query(listTransactionQuery, userId)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer listTransactionResult.Close()

	var listUserTransaction []dto.User

	for listTransactionResult.Next(){
		var userTransaction dto.User
		err = listTransactionResult.Scan(&userTransaction.UserId, &userTransaction.Description, &userTransaction.Value, &userTransaction.LastUpdatedAt)
		if err != nil {
			return err
		}
		// Format CreatedAt
		dt,_ := time.Parse(time.RFC3339, userTransaction.LastUpdatedAt)
		userTransaction.LastUpdatedAt = dt.Format(constant.DateTimeLayout)

		listUserTransaction = append(listUserTransaction, userTransaction)
	}

	err = service.Cache.SetWithError(constant.RedisPrefixKeyAllTransaction + userId, listUserTransaction, 0)
	if err != nil {
		logger.Error("Error update redis", err.Error())
		return err
	}

	return nil
}

/*
	Update bought mobile card
 */
func (service *RedisService) UpdateBoughtMobileCardRedis(userId string) error {
	// 2. Update bought mobile cards
	var mobileCardSuccessfully 	[]dto.MobileCard

	getAllBoughtMobileCardQuery := `SELECT uuid_from_bin(mobile_card.Id), mobile_card_vendor.Name, mobile_card.VendorCode, mobile_card.Serial, mobile_card.Code, mobile_card.Value, game_mobile_card.CreatedAt 
									FROM game_mobile_card, mobile_card, mobile_card_vendor 
									WHERE mobile_card_vendor.VendorCode = mobile_card.VendorCode AND game_mobile_card.MobileCardId = mobile_card.Id AND uuid_from_bin(game_mobile_card.UserId) = ?
									ORDER BY game_mobile_card.CreatedAt DESC;`
	getAllBoughtMobileCardResult, err := service.MySql.DbContext.Query(getAllBoughtMobileCardQuery, userId)
	if err != nil {
		logger.Error(err.Error())
		return err
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
			return err
		}
		mobileCardSuccessfully = append(mobileCardSuccessfully, mobileCard)
	}

	err = service.Cache.SetWithError(constant.RedisPrefixKeyBoughtMobileCards + userId, mobileCardSuccessfully, 0)
	if err != nil {
		logger.Error("Error update redis", err.Error())
		return err
	}

	return nil
}

/*
	Update all prize
 */
func (service *RedisService) UpdateAllPrizeRedis() error{

	prizeQuery 	:= `SELECT uuid_from_bin(Id) AS Id, Name, Value, Description, CreatedAt, LastUpdatedAt 
								FROM user_prize 
								ORDER BY user_prize.LastUpdatedAt DESC;`

	prizeResult, err := service.MySql.DbContext.Query(prizeQuery)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer prizeResult.Close()

	var listPrize []dto.Prize

	for prizeResult.Next(){
		var prize dto.Prize
		err = prizeResult.Scan(&prize.Id, &prize.Name, &prize.Value, &prize.Description, &prize.CreatedAt, &prize.LastUpdatedAt)
		if err != nil {
			return err
		}

		listPrize = append(listPrize, prize)
	}

	err = service.Cache.SetWithError(constant.RedisPrefixKeyAllPrize, listPrize, 0)
	if err != nil {
		logger.Error("Error update redis", err.Error())
		return err
	}

	return nil
}


/*
	Update Vendor
 */

func (service *RedisService) UpdateAllVendorRedis() error{
	vendorQuery 	:= `SELECT uuid_from_bin(Id) AS Id, Name, VendorCode, Status, CreatedAt, LastUpdatedAt 
						FROM mobile_card_vendor 
						ORDER BY LastUpdatedAt DESC;`

	vendorResult, err := service.MySql.DbContext.Query(vendorQuery)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer vendorResult.Close()

	var listVendor []dto.MobileCardVendor

	for vendorResult.Next(){
		var vendor dto.MobileCardVendor
		err = vendorResult.Scan(&vendor.Id, &vendor.Name, &vendor.VendorCode, &vendor.Status, &vendor.CreatedAt, &vendor.LastUpdatedAt)
		if err != nil {
			return err
		}

		listVendor = append(listVendor, vendor)
	}

	err = service.Cache.SetWithError(constant.RedisPrefixKeyAllVendor, listVendor, 0)
	if err != nil {
		logger.Error("Error update redis", err.Error())
		return err
	}

	return nil
}


/*
	Update wallet by Id
 */
func (service *RedisService) UpdateUserWalletRedis(userId string) error{

	// Get inviting status
	userInvitingQuery := `SELECT uuid_from_bin(Id)
				FROM game_inviting
				WHERE uuid_from_bin(InvitedUser) = ?;`
	userInvitingResult, err := service.MySql.DbContext.Query(userInvitingQuery, userId)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer userInvitingResult.Close()

	var isInvited bool
	if userInvitingResult.Next(){
		isInvited = true
	}else{
		isInvited = false
	}

	// Get Wallet
	userWalletQuery := `SELECT IFNULL(SUM(user_wallet.Value), 0) AS Wallet  
				FROM user_wallet 
				WHERE uuid_from_bin(UserId) = ?;`
	userWalletResult, err := service.MySql.DbContext.Query(userWalletQuery, userId)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}

	defer userWalletResult.Close()

	var wallet int
	if userWalletResult.Next(){
		err = userWalletResult.Scan(&wallet)
		if err != nil {
			return  err
		}
	}

	user := dto.User{
		UserId:	userId,
		Wallet:	wallet,
		IsInvited: isInvited,
	}

	err = service.Cache.SetWithError(constant.RedisPrefixKeyUserWallet + userId, user, 0)
	if err != nil {
		logger.Error("Error update redis", err.Error())
		return err
	}
	return nil
}

/*
	Update wallet by Id
*/
func (service *RedisService) ResetRedis() error{

	usersQuery := `SELECT uuid_from_bin(Id) FROM sso_user`
	usersResult, err := service.MySql.DbContext.Query(usersQuery)
	if err != nil {
		service.MySql.HandleError(err)
		return err
	}
	defer usersResult.Close()

	for usersResult.Next(){
		var userId string
		err = usersResult.Scan(&userId)
		if err != nil {
			return  err
		}
		affectedWallet := service.Cache.DeleteItem(constant.RedisPrefixKeyUserWallet + userId)
		fmt.Printf("Wallet %s: %d\n", userId, affectedWallet)
		affectedBoughtCards := service.Cache.DeleteItem(constant.RedisPrefixKeyBoughtMobileCards + userId)
		fmt.Printf("BoughtCard %s: %d\n", userId, affectedBoughtCards)
		affectedTransaction := service.Cache.DeleteItem(constant.RedisPrefixKeyAllTransaction + userId)
		fmt.Printf("Transaction %s: %d\n", userId, affectedTransaction)
	}
	affectedMobileCardVendor := service.Cache.DeleteItem(constant.RedisPrefixKeyAllVendor)
	fmt.Printf("MobileCardVendors: %d\n", affectedMobileCardVendor)
	affectedPrize := service.Cache.DeleteItem(constant.RedisPrefixKeyAllPrize)
	fmt.Printf("Prizes: %d\n", affectedPrize)

	return nil
}
