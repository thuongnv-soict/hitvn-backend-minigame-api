package controller

import (
	"context"
	"g-tech.com/dto"
	"g-tech.com/gerror"
	"g-tech.com/infrastructure/controller"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/response"
	"g-tech.com/infrastructure/util"
	"g-tech.com/module/minigame/service"
	"github.com/labstack/echo"
)

type LotteryController struct {
	controller.BaseController
	Service     service.ILotteryService
}

func NewLotteryController(lotteryService service.ILotteryService) *LotteryController{
	return &LotteryController{
		Service: lotteryService,
	}
}

/*
	Update mobile card
*/
func (controller *LotteryController) GetSelectedNumbers(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	userId := echo.Param("userId")

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// 3. Retrieve data
	result, err := controller.Service.GetSelectedNumbers(ctx, userId)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}
	return controller.WriteSuccess(echo, result)
}

/*
	Create new Mobile Card
*/
func (controller *LotteryController) CreateLotteryNumber(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	lotteryPlayer := dto.LotteryPlayer{}
	err := echo.Bind(&lotteryPlayer)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&lotteryPlayer); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	errorCode, err := controller.Service.CreateLotteryNumber(ctx, lotteryPlayer)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	if errorCode != 0 {
		message, errRes := response.NewErrorResponse(errorCode, "", util.FuncName())
		return controller.WriteStatusConflict(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

