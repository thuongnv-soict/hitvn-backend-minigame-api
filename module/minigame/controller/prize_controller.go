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

type PrizeController struct {
	controller.BaseController
	Service     service.IPrizeService
}

func NewPrizeController(prizeService service.IPrizeService) *PrizeController{
	return &PrizeController{
		Service: prizeService,
	}
}


/*
	Create new prize
 */
func (controller *PrizeController) CreatePrize(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	prize := dto.Prize{}
	err := echo.Bind(&prize)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&prize); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.CreatePrize(ctx, prize)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Get all prize
 */
func (controller *PrizeController) GetAllPrize(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listPrize, err := controller.Service.GetAllPrize(ctx, pageSize, pageIndex)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listPrize)
}

/*
	Update prize
*/
func (controller *PrizeController) UpdatePrize(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	prize := dto.Prize{}
	err := echo.Bind(&prize)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&prize); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.UpdatePrize(ctx, prize)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Delete prize
*/
func (controller *PrizeController) DeletePrize(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	prizeId := echo.Param("prizeId")

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := controller.Service.DeletePrize(ctx, prizeId)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}



