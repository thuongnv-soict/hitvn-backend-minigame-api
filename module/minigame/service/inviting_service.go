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
	"github.com/pkg/errors"
	"time"
)

type IInvitingService interface {
	CreateNewInvitation(ctx context.Context, invitation dto.Invitation) (int, int, error)
	GetInvitingCode(ctx context.Context, phoneNumber string) (string, error)
	GenerateInvitingCode(ctx context.Context, phoneNumber string) (string, error)
}

type InvitingService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	RedisService	RedisService
	ConfigService 	ConfigService
	Timeout    		time.Duration
}

func NewInvitingService (dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, configService ConfigService, timeout time.Duration) IInvitingService {
	service := InvitingService{}
	service.Cache = cache
	service.RedisService = redisService
	service.ConfigService = configService
	service.MySql.SetDbContext(dbContext)
	service.Timeout = timeout
	return &service
}

/*
	Get invited code
	If it is not exited, generate it from phone number
 */
func (service *InvitingService) GetInvitingCode(ctx context.Context, phoneNumber string) (string, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	getCodeQuery := `SELECT IFNULL(Code, "") FROM sso_user WHERE PhoneNumber = ?;`
	getCodeResult, err := service.MySql.DbContext.Query(getCodeQuery, phoneNumber)
	if err != nil {
		service.MySql.HandleError(err)
		return "", err
	}
	defer getCodeResult.Close()

	var code string
	if getCodeResult.Next(){
		err := getCodeResult.Scan(&code)
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}
	} else {
		return "", errors.New("Cannot find phone number")
	}

	// Code existed
	if code != "" {
		return code, nil
	}

	//	Code is not existed
	code, err = service.GenerateInvitingCode(ctx, phoneNumber)
	if err != nil {
		return "", err
	}

	return code, nil
}

/*
	Generate minigame code
 */
func (service *InvitingService) GenerateInvitingCode(ctx context.Context, phoneNumber string) (string, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	code, err := util.EncodeInvitedCode(phoneNumber)
	if err != nil {
		return "", err
	}
	if code == "" {
		return "", errors.New("Cannot encode phone number")
	}

	updateCodeStatement := `UPDATE sso_user SET Code = ? WHERE PhoneNumber = ?`
	updateCodeResult, err := service.MySql.DbContext.Exec(updateCodeStatement, code, phoneNumber)
	if err != nil {
		service.MySql.HandleError(err)
		return "", err
	}

	rowsAffected, err := updateCodeResult.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0{
		return "", errors.New("Cannot update code")
	}

	return code, nil
}


/*
	Check InvitedCode inserted (already inserted)
 */
func (service *InvitingService) checkStatusCodeInserted(invitedUserId string) (bool, error){
	invitingUserQuery := `SELECT uuid_from_bin(InvitingUser) FROM game_inviting WHERE uuid_from_bin(InvitedUser) = ?;`
	invitingUserResult, err := service.MySql.DbContext.Query(invitingUserQuery, invitedUserId)
	if err != nil {
		service.MySql.HandleError(err)
		return false, err
	}
	defer invitingUserResult.Close()

	if invitingUserResult.Next(){
		return false, nil
	}

	return true, nil
}

/*
 	Get InvitingUserId from invitedCode
 */
func (service *InvitingService) getInvitingUserID(invitedUserId string, invitedCode string) (string, bool, error){
	invitingUserId := ""

	if invitedUserId == "" {
		return invitingUserId, false, errors.New("Don't insert empty code")
	}

	invitingUserQuery := `SELECT uuid_from_bin(Id) FROM sso_user WHERE Code = ?;`
	invitingUserResult, err := service.MySql.DbContext.Query(invitingUserQuery, invitedCode)
	if err != nil {
		service.MySql.HandleError(err)
		return invitingUserId, false, err
	}
	defer invitingUserResult.Close()

	if invitingUserResult.Next(){
		err := invitingUserResult.Scan(&invitingUserId)
		if err != nil {
			logger.Error(err.Error())
			return invitingUserId, false, err
		}
		return invitingUserId, true, nil
	}

	return invitingUserId, false, nil
}

///*
//	Get Prize
// */
//func (summary *InvitingService) GetPrize(prizeName string) (dto.Prize, bool, error){
//	// Check if it exist in Redis
//	prize := dto.Prize{}
//	prize, status, err := summary.GetPrizeRedis(prizeName)
//	if err == redis.Nil {
//		err = summary.RedisService.UpdateAllPrizeRedis()
//		if err == nil {
//			prize, status, err = summary.GetPrizeRedis(prizeName)
//		}
//	}
//
//	if err == nil {
//		return prize, status, nil
//	}
//
//	// Update error  or get error
//	prize, status, err = summary.GetPrizeSQL(prizeName)
//	return prize, status, err
//}
//
///*
//	Get prize redis
//*/
//func (summary *InvitingService) GetPrizeRedis(prizeName string) (dto.Prize, bool, error){
//	var prize dto.Prize
//	var listPrize []dto.Prize
//
//	result, err := summary.Cache.GetWithError(constant.RedisPrefixKeyAllPrize)
//	if err == redis.Nil {
//		return prize, false, nil
//	} else if err != nil {
//		logger.Error(err.Error())
//		return prize, false, err
//	}
//
//	err = json.Unmarshal([]byte(result), &listPrize)
//	if err != nil {
//		logger.Error(err.Error())
//		return prize, false, err
//	}
//	for _, p := range listPrize {
//		if p.Name == prizeName {
//			prize = p
//			break
//		}
//	}
//	return prize, true, nil
//}
//
///*
//	Get prize SQL
// */
//func (summary *InvitingService) GetPrizeSQL(prizeName string) (dto.Prize, bool, error){
//	prize := dto.Prize{}
//
//	invitingUserQuery := `SELECT uuid_from_bin(Id), Value, Description FROM user_prize WHERE Name = ?;`
//	invitingUserResult, err := summary.MySql.DbContext.Query(invitingUserQuery, prizeName)
//	if err != nil {
//		summary.MySql.HandleError(err)
//		return prize, false, err
//	}
//	defer invitingUserResult.Close()
//
//	if invitingUserResult.Next(){
//		err := invitingUserResult.Scan(&prize.Id, &prize.Value, &prize.Description)
//		if err != nil {
//			logger.Error(err.Error())
//			return prize, false, err
//		}
//		return prize, true, nil
//	}
//
//	return prize, false, nil
//}


/*
	Create new invitation
 */
func (service *InvitingService) CreateNewInvitation(ctx context.Context, invitation dto.Invitation) (int, int, error){
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	/*
		Check requirement
	 */
	//	Check invited code inserted
	status, err := service.checkStatusCodeInserted(invitation.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	} else{
		if status == false {
			return 0, gerror.ErrorUserHasInsertedCode, nil
		}
	}

	// Check invited code valid
	invitingUserId, status, err := service.getInvitingUserID(invitation.UserId, invitation.Code)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	} else {
		if status == false {
			return 0, gerror.ErrorInvitingUserDoesNotExisted, nil
		}else{
			if invitingUserId == invitation.UserId{
				return 0, gerror.ErrorInvitingUserIsYou,  nil
			}
		}
	}

	// Get prize for Inviting User
	invitingPrize, status, err := service.ConfigService.GetPrize(constant.ProgramInvitingName)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	} else {
		if status == false {
			return 0, gerror.ErrorInvitingProgramNotFound, nil
		}
	}

	// Get prize for InvitedUser
	invitedPrize, status, err := service.ConfigService.GetPrize(constant.ProgramInvitedName)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	} else {
		if status == false {
			return 0, gerror.ErrorInvitedProgramNotFound, nil
		}
	}

	/*
		Apply invited code
	 */
	// Start transaction
	tx, err := service.MySql.DbContext.Begin()
	
	//	Insert statistic (game_inviting)
	walletId, status, err := service.CreateInvitingRecord(tx, invitingUserId, invitation.UserId)
	if err != nil{
		logger.Error(err.Error())
		_ = tx.Rollback()
		return 0, 0, err
	}


	//	Insert History for minigame user(table: User_wallet)
	status, err = service.CreateNewHistory(tx, invitingUserId, walletId, invitingPrize.Id, invitingPrize.Value)
	if err != nil{
		logger.Error(err.Error())
		_ = tx.Rollback()
		return 0, 0, err
	}

	//	Insert History for minigame user(table: User_wallet)
	status, err = service.CreateNewHistory(tx, invitation.UserId, walletId, invitedPrize.Id, invitedPrize.Value)
	if err != nil{
		logger.Error(err.Error())
		_ = tx.Rollback()
		return 0, 0, err
	}

	_ = tx.Commit()

	/*
	 	Update Redis
	 */
	// For minigame user
	err = service.RedisService.UpdateTransactionRedis(invitingUserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	}

	// For invited user
	err = service.RedisService.UpdateTransactionRedis(invitation.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	}

	//	Update inviting user wallet
	err = service.RedisService.UpdateUserWalletRedis(invitingUserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	}

	//	Update invited user wallet
	err = service.RedisService.UpdateUserWalletRedis(invitation.UserId)
	if err != nil {
		logger.Error(err.Error())
		return 0, 0, err
	}

	return invitedPrize.Value, 0, nil
}


/*
	Insert Inviting Record (table: game_inviting)
 */
func (service *InvitingService) CreateInvitingRecord(tx *sql.Tx, invitingUserId string, invitedUserId string) (string, bool, error){
	insertInvitingRecordStatement := `INSERT INTO game_inviting(Id, InvitingUser, InvitedUser, WalletId) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`

	walletId := util.NewUuid()
	_, err := tx.Exec(insertInvitingRecordStatement, util.NewUuid(), invitingUserId, invitedUserId, walletId)
	if err != nil {
		fmt.Println(err.Error())
		return walletId, false, err
	}

	return walletId, true, nil
}


/*
	Create History
 */
func (service *InvitingService) CreateNewHistory(tx *sql.Tx, userId string, walletId string, prizeId string, value int) (bool, error) {

	postInvitingRecordStatement := `INSERT INTO user_wallet(Id, UserId, PrizeId, Value) VALUES (uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?), ?);`

	_, err := tx.Exec(postInvitingRecordStatement, walletId, userId, prizeId, value)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	return true, nil
}
