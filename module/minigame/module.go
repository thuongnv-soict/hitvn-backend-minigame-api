package minigame

import (
	"database/sql"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/module/minigame/controller"
	"g-tech.com/module/minigame/service"
	"github.com/labstack/echo"
	"time"
)

var invitingController 			*controller.InvitingController
var prizeController 			*controller.PrizeController
var mobileCardController 		*controller.MobileCardController
var mobileCardVendorController 	*controller.MobileCardVendorController
var userController 				*controller.UserController
var readDailyController 		*controller.ReadDailyController
var lotteryController			*controller.LotteryController

func Initialize(e *echo.Echo, dbContext *sql.DB, cache cache.CacheManager, timeout time.Duration){
	redisService 				:= service.NewRedisService(dbContext, cache, timeout)
	configService 				:= service.NewConfigService(dbContext, cache, redisService, timeout)

	prizeService 				:= service.NewPrizeService(dbContext, cache, redisService, timeout)
	prizeController 			= controller.NewPrizeController(prizeService)

	invitingService 			:= service.NewInvitingService(dbContext, cache, redisService, configService, timeout)
	invitingController	 		= controller.NewInvitingController(invitingService)

	lotteryService 				:= service.NewLotteryService(dbContext, cache, redisService, configService, timeout)
	lotteryController	 		= controller.NewLotteryController(lotteryService)


	userService 				:= service.NewUserService(dbContext, cache, redisService, timeout)
	userController 				= controller.NewUserController(userService)

	readDailyService 			:= service.NewReadDailyService(dbContext, cache, redisService, timeout)
	readDailyController 		= controller.NewReadDailyController(readDailyService)

	mobileCardService 			:= service.NewMobileCardService(dbContext, cache, redisService, configService, timeout)
	mobileCardController 		= controller.NewMobileCardController(mobileCardService)

	mobileCardVendorService 		:= service.NewMobileCardVendorService(dbContext, cache, redisService, timeout)
	mobileCardVendorController 		= controller.NewMobileCardVendorController(mobileCardVendorService)

	initRouter(e)
}

func initRouter(e *echo.Echo){

	e.GET("/game/api/v1.0/mini-game/statistic/wallet/user/:userId", userController.GetUserWallet)


	// Read daily
	e.POST("/game/api/v1.0/mini-game/read-daily", readDailyController.CreateNewReadDaily)

	// Invitation
	e.POST("/game/api/v1.0/mini-game/invitation", invitingController.CreateNewInvitation)
	e.GET("/game/api/v1.0/mini-game/invitation/code/:phoneNumber", invitingController.GetInvitingCode)

	// MobileCard
	e.POST("/game/api/v1.0/mini-game/exchange-mobile-card/exchange", mobileCardController.ExchangeMobileCard)
	e.GET("/game/api/v1.0/mini-game/exchange-mobile-card/list/bought/:userId", mobileCardController.GetListBoughtMobileCard)

	// Prize Management
	e.GET("/game/api/v1.0/prize-management/list", prizeController.GetAllPrize)
	e.POST("/game/api/v1.0/prize-management/add", prizeController.CreatePrize)
	e.PUT("/game/api/v1.0/prize-management/update", prizeController.UpdatePrize)
	e.DELETE("/game/api/v1.0/prize-management/delete/:prizeId", prizeController.DeletePrize)


	/*
		Mobile Card Vendor
	 */
	//	For app
	e.GET("/game/api/v1.0/mini-game/exchange-mobile-card/list/vendor", mobileCardVendorController.GetListActiveVendor)
	e.GET("/game/api/v1.0/mobile-card-vendor/statistic/active-mobile-card/:vendor", mobileCardVendorController.GetListQuantityActiveMobileCard)
	//	For web
	e.GET("/game/api/v1.0/mobile-card-vendor/list/active", mobileCardVendorController.GetListActiveVendor)
	e.GET("/game/api/v1.0/mobile-card-vendor/list/all", mobileCardVendorController.GetListAllVendor)
	e.POST("/game/api/v1.0/mobile-card-vendor/add", mobileCardVendorController.CreateMobileCardVendor)
	e.PUT("/game/api/v1.0/mobile-card-vendor/update", mobileCardVendorController.UpdateMobileCardVendor)
	e.DELETE("/game/api/v1.0/mobile-card-vendor/delete/:mobileCardVendorId", mobileCardVendorController.DeleteMobileCardVendor)

	/*
		Mobile Card Management
	 */
	e.GET("/game/api/v1.0/mobile-card/list", mobileCardController.GetMobileCard)
	e.POST("/game/api/v1.0/mobile-card/add", mobileCardController.CreateMobileCard)
	e.PUT("/game/api/v1.0/mobile-card/update", mobileCardController.UpdateMobileCard)
	e.DELETE("/game/api/v1.0/mobile-card/delete/:mobileCardId", mobileCardController.DeleteMobileCard)

	/*
		Lottery
	 */
	e.GET("/game/api/v1.0/lottery/selected-numbers/:userId", lotteryController.GetSelectedNumbers)
	e.POST("/game/api/v1.0/lottery/add", lotteryController.CreateLotteryNumber)

	// Transaction
	e.GET("/game/api/v1.0/mini-game/statistic/transaction/list/:userId", userController.ListTransactions)
	e.GET("/game/api/v1.0/mini-game/statistic/transaction/received/:userId", userController.GetReceivedTransaction)
	e.GET("/game/api/v1.0/mini-game/statistic/transaction/used/:userId", userController.GetUsedTransaction)

	//  Redis
	//e.DELETE("/game/api/v1.0/reset-redis", userController.ResetRedis)
}