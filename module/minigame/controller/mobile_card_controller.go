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

type MobileCardController struct {
	controller.BaseController
	Service     service.IMobileCardService
}

func NewMobileCardController(mobileCardService service.IMobileCardService) *MobileCardController{
	return &MobileCardController{
		Service: mobileCardService,
	}
}
/*
	Get all mobile card
 */
/*
	Update mobile card
*/
func (controller *MobileCardController) GetMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))


	mobileCardFilter := dto.MobileCardFilter{}
	//mobileCardFilter.Name = echo.QueryParam("Name")
	//mobileCardFilter.Value, _ = strconv.Atoi(echo.QueryParam("Value"))
	//mobileCardFilter.Status, _ = strconv.Atoi(echo.QueryParam("Status"))
	err := echo.Bind(&mobileCardFilter)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&mobileCardFilter); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listMobileCard, err := controller.Service.GetMobileCard(ctx, mobileCardFilter, pageSize, pageIndex)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listMobileCard)
}

/*
	Create new Mobile Card
 */
func (controller *MobileCardController) CreateMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCard := dto.MobileCard{}
	err := echo.Bind(&mobileCard)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&mobileCard); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.CreateMobileCard(ctx, mobileCard)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Update mobile card
 */
func (controller *MobileCardController) UpdateMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCard := dto.MobileCard{}
	err := echo.Bind(&mobileCard)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&mobileCard); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.UpdateMobileCard(ctx, mobileCard)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Delete mobile card
*/
func (controller *MobileCardController) DeleteMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCardId := echo.Param("mobileCardId")

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := controller.Service.DeleteMobileCard(ctx, mobileCardId)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Exchange mobile card
 */
func (controller *MobileCardController) ExchangeMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	userExchange := dto.UserExchange{}
	err := echo.Bind(&userExchange)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&userExchange); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	mobileCard, errorCode, err := controller.Service.ExchangeMobileCard(ctx, userExchange)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	if errorCode != 0 {
		message, errRes := response.NewErrorResponse(errorCode, "", util.FuncName())
		return controller.WriteStatusConflict(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, mobileCard)
}


/*
	Get list bought mobile card
*/
func (controller *MobileCardController) GetListBoughtMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	userId 			:= echo.Param("userId")
	pageSize, _ 	:= strconv.Atoi(echo.QueryParam("pageSize"))
	pageIndex, _ 	:= strconv.Atoi(echo.QueryParam("pageIndex"))

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listBoughtMobileCard, err := controller.Service.GetListBoughtMobileCard(ctx, userId, pageSize, pageIndex)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listBoughtMobileCard)
}