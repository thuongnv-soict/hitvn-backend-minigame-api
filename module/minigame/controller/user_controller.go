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
	"strconv"
)

type UserController struct {
	controller.BaseController
	Service     service.IUserService
}

func NewUserController(UserService service.IUserService) *UserController{
	return &UserController{
		Service: UserService,
	}
}

/*
	Get User Wallet
	+ wallet
	+ current day of week
	+ is invited
	+ read days of week
 */
func (controller *UserController) GetUserWallet(echo echo.Context) error {
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	userId := echo.Param("userId")

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// 3. Retrieve data
	result, err := controller.Service.GetWalletByUserId(ctx, userId)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}
	return controller.WriteSuccess(echo, result)
}

/*
	Get all transaction
 */
func (controller *UserController) ListTransactions(echo echo.Context) error{
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	userId 			:= echo.Param("userId")
	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var listTransactions []dto.User
	listTransactions, err := controller.Service.ListTransactions(ctx, userId, 0, pageSize, pageIndex)
	if err != nil {
			message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
			return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listTransactions)
}

/*
	Get received transaction
*/
func (controller *UserController) GetReceivedTransaction(echo echo.Context) error{
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	userId 			:= echo.Param("userId")
	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var listTransactions []dto.User
	listTransactions, err := controller.Service.ListTransactions(ctx, userId, 1, pageSize, pageIndex)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listTransactions)
}

/*
	Get used transaction
*/
func (controller *UserController) GetUsedTransaction(echo echo.Context) error{
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	userId := echo.Param("userId")
	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var listTransactions []dto.User
	listTransactions, err := controller.Service.ListTransactions(ctx, userId, -1, pageSize, pageIndex)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listTransactions)
}

/*
	Get used transaction
*/
func (controller *UserController) ResetRedis(echo echo.Context) error{
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var listTransactions []dto.User
	err := controller.Service.ResetRedis(ctx)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listTransactions)
}

