package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"github.com/go-redis/redis"
	"time"
)

type IMobileCardVendorService interface {
	GetListAllVendor(ctx context.Context) ([]dto.MobileCardVendor, error)
	GetListActiveVendor(ctx context.Context) ([]dto.MobileCardVendor, error)
	GetListQuantityActiveMobileCard(ctx context.Context, vendorName string) ([]dto.MobileCardVendor, error)
	CreateMobileCardVendor(ctx context.Context, mobileCardVendor dto.MobileCardVendor) error
	UpdateMobileCardVendor(ctx context.Context, vendor dto.MobileCardVendor) error
	DeleteMobileCardVendor(ctx context.Context, vendorId string) error

}

type MobileCardVendorService struct {
	MySql 			repository.MySqlRepository
	Cache 			cache.CacheManager
	RedisService 	RedisService
	Timeout    		time.Duration
}


func NewMobileCardVendorService(dbContext *sql.DB, cache cache.CacheManager, redisService RedisService, timeout time.Duration) IMobileCardVendorService {
	service := MobileCardVendorService{}
	service.MySql.SetDbContext(dbContext)
	service.Cache = cache
	service.RedisService = redisService
	service.Timeout = timeout
	return &service
}

/*
	Get list quantity active mobile card
 */
func (service *MobileCardVendorService) GetListQuantityActiveMobileCard(ctx context.Context, vendorName string) ([]dto.MobileCardVendor, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	var listMobileCard []dto.MobileCardVendor
	getListActiveMobileCard := `SELECT mobile_card_vendor.Name, mobile_card.Value, mobile_card.Status, COUNT(mobile_card.Id) AS Quantity
								FROM mobile_card, mobile_card_vendor
								WHERE mobile_card.VendorCode = mobile_card_vendor.Name AND mobile_card_vendor.Name = ? AND mobile_card_vendor.Status = ? AND mobile_card.Status = ?
								GROUP BY mobile_card.Value;`

	getListActiveMobileCardResult, err := service.MySql.DbContext.Query(getListActiveMobileCard, vendorName, constant.StatusMobileCardVendorActive, constant.StatusMobileCardReady)
	if err != nil {
		logger.Error(err.Error())
		return listMobileCard, err
	}
	defer getListActiveMobileCardResult.Close()

	for getListActiveMobileCardResult.Next(){
		var mobileCard dto.MobileCardVendor
		err := getListActiveMobileCardResult.Scan(&mobileCard.Name, &mobileCard.Value, &mobileCard.Status, &mobileCard.Quantity)
		if err != nil{
			return listMobileCard, err
		}
		listMobileCard = append(listMobileCard, mobileCard)
	}

	return listMobileCard, nil
}

/*
	Get list vendor
*/
func (service *MobileCardVendorService) GetListAllVendor(ctx context.Context) ([]dto.MobileCardVendor, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	// Check if it exist in Redis
	var listVendor []dto.MobileCardVendor
	listVendor, err := service.GetListVendorRedis(ctx)
	if err == redis.Nil {
		err = service.RedisService.UpdateAllVendorRedis()
		if err == nil {
			listVendor, err = service.GetListVendorRedis(ctx)
		}
	}

	if err == nil {
		return listVendor, nil
	}

	// Update error  or get error
	listVendor, err = service.GetListVendorSQL(ctx)
	return listVendor, err

}

/*
	Get list vendor Redis
 */
func (service *MobileCardVendorService) GetListVendorRedis(ctx context.Context) ([]dto.MobileCardVendor, error) {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	result, err := service.Cache.GetWithError(constant.RedisPrefixKeyAllVendor)
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var allVendor []dto.MobileCardVendor
	err = json.Unmarshal([]byte(result), &allVendor)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return allVendor, nil
}

/*
	Get list vendor SQL
*/
func (service *MobileCardVendorService) GetListVendorSQL(ctx context.Context) ([]dto.MobileCardVendor, error) {
	vendorQuery 	:= `SELECT uuid_from_bin(Id) AS Id, Name, VendorCode, Status, CreatedAt, LastUpdatedAt 
								FROM mobile_card_vendor 
								ORDER BY LastUpdatedAt DESC;`

	vendorResult, err := service.MySql.DbContext.Query(vendorQuery)
	if err != nil {
		service.MySql.HandleError(err)
		return nil, err
	}
	defer vendorResult.Close()

	var listVendor []dto.MobileCardVendor

	for vendorResult.Next(){
		var vendor dto.MobileCardVendor
		err = vendorResult.Scan(&vendor.Id, &vendor.Name, &vendor.VendorCode, &vendor.Status, &vendor.CreatedAt, &vendor.LastUpdatedAt)
		if err != nil {
			return nil, err
		}

		listVendor = append(listVendor, vendor)
	}

	return listVendor, nil
}

/*
	Get list active vendor
*/
func (service *MobileCardVendorService) GetListActiveVendor(ctx context.Context) ([]dto.MobileCardVendor, error) {
	allVendor, err := service.GetListAllVendor(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var listActiveVendor []dto.MobileCardVendor
	for _, vendor := range allVendor {
		if vendor.Status == constant.StatusMobileCardVendorActive{
			listActiveVendor = append(listActiveVendor, vendor)
		}
	}

	return listActiveVendor, nil
}

/*
	Create mobile card vendor
 */

func (service *MobileCardVendorService) CreateMobileCardVendor(ctx context.Context, mobileCardVendor dto.MobileCardVendor) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	createMobileCardVendorStatement, err := service.MySql.DbContext.Prepare(`INSERT INTO mobile_card_vendor(Id, Name, VendorCode, Status) VALUES (uuid_to_bin(?), ?, ?, ?);`)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer createMobileCardVendorStatement.Close()


	createMobileCardVendorResult, err := createMobileCardVendorStatement.Exec(util.NewUuid(), mobileCardVendor.Name, mobileCardVendor.VendorCode, mobileCardVendor.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	rowsAffected, err := createMobileCardVendorResult.RowsAffected()
	if err != nil {
		return errors.New("Cannot get row affected")
	}
	if rowsAffected == 0 {
		return errors.New("No row affected")
	}

	err = service.RedisService.UpdateAllVendorRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


/*
	Update Mobile Card
*/
func (service *MobileCardVendorService) UpdateMobileCardVendor(ctx context.Context, mobileCardVendor dto.MobileCardVendor) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	updateMobileCardStatement, err := service.MySql.DbContext.Prepare(`UPDATE mobile_card_vendor
														SET Name = ?, 
															VendorCode = ?,
															Status = ?
														WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	_, err = updateMobileCardStatement.Exec(mobileCardVendor.Name, mobileCardVendor.VendorCode, mobileCardVendor.Status, mobileCardVendor.Id)
	if err != nil {
		return err
	}
	defer updateMobileCardStatement.Close()

	err = service.RedisService.UpdateAllVendorRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

/*
	Delete Mobile Card Vendor
*/
func (service *MobileCardVendorService) DeleteMobileCardVendor(ctx context.Context, vendorId string) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, service.Timeout)
	defer cancel()

	deleteVendorStatement, err := service.MySql.DbContext.Prepare(`DELETE FROM mobile_card_vendor
																		WHERE uuid_from_bin(Id) = ?;`)
	if err != nil {
		return err
	}

	deleteVendorResult, err := deleteVendorStatement.Exec(vendorId)
	if err != nil {
		return err
	}
	defer deleteVendorStatement.Close()

	rowsAffected, err := deleteVendorResult.RowsAffected()
	if err != nil {
		return errors.New("Cannot get row affected")
	}
	if rowsAffected == 0 {
		return errors.New("No row deleted")
	}

	err = service.RedisService.UpdateAllVendorRedis()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}