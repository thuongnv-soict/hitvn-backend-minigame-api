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

type MobileCardVendorController struct {
	controller.BaseController
	Service     service.IMobileCardVendorService
}

func NewMobileCardVendorController(mobileCardVendorService service.IMobileCardVendorService) *MobileCardVendorController{
	return &MobileCardVendorController{
		Service: mobileCardVendorService,
	}
}

/*
	Get list active vendor
*/
func (controller *MobileCardVendorController) GetListActiveVendor(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listActiveVendor, err := controller.Service.GetListActiveVendor(ctx)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listActiveVendor)
}

/*
	Get list vendor
*/
func (controller *MobileCardVendorController) GetListAllVendor(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listVendor, err := controller.Service.GetListAllVendor(ctx)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listVendor)
}

/*
	Get list quantity active mobile card
*/
func (controller *MobileCardVendorController) GetListQuantityActiveMobileCard(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	vendorName := echo.Param("vendor")

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listQuantityMobileCard, err := controller.Service.GetListQuantityActiveMobileCard(ctx, vendorName)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorRetrieveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, listQuantityMobileCard)
}

func (controller *MobileCardVendorController) CreateMobileCardVendor(echo echo.Context) error {
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCardVendor := dto.MobileCardVendor{}
	err := echo.Bind(&mobileCardVendor)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&mobileCardVendor); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.CreateMobileCardVendor(ctx, mobileCardVendor)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}

/*
	Update mobile card vendor
*/
func (controller *MobileCardVendorController) UpdateMobileCardVendor(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCardVendor := dto.MobileCardVendor{}
	err := echo.Bind(&mobileCardVendor)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&mobileCardVendor); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = controller.Service.UpdateMobileCardVendor(ctx, mobileCardVendor)

	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}


/*
	Delete mobile card vendor
*/
func (controller *MobileCardVendorController) DeleteMobileCardVendor(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	mobileCardVendorId := echo.Param("mobileCardVendorId")
	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := controller.Service.DeleteMobileCardVendor(ctx, mobileCardVendorId)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	return controller.WriteSuccessEmptyContent(echo)
}
