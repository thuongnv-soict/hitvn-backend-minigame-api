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

type InvitingController struct {
	controller.BaseController
	Service service.IInvitingService
}

func NewInvitingController(invitingService service.IInvitingService) *InvitingController{
	return &InvitingController{
		Service: invitingService,
	}
}

/*
	Create new invitation
 */
func (controller *InvitingController) CreateNewInvitation(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	invitation := dto.Invitation{}
	err := echo.Bind(&invitation)
	if err != nil {
		message, errorRes := response.NewErrorResponse(gerror.ErrorBindData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errorRes)
	}

	// 3. validate object
	if ok, err := controller.IsValid(&invitation); !ok && err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorValidData, err.Error(), util.FuncName())
		return controller.WriteBadRequest(echo, message, errRes)
	}

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	prizeValue, errorCode, err := controller.Service.CreateNewInvitation(ctx, invitation)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorSaveData, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}
	if errorCode != 0 {
		message, errRes := response.NewErrorResponse(errorCode, "", util.FuncName())
		return controller.WriteStatusConflict(echo, message, errRes)
	}

	return controller.WriteSuccess(echo, dto.User{
		UserId: invitation.UserId,
		IsInvited: true,
		Value: prizeValue,
	})
}

/*
	Get invited code
*/
func (controller *InvitingController) GetInvitingCode(echo echo.Context) error{
	// 0. log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	phoneNumber := echo.Param("phoneNumber")

	// 4. Define Context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	invitingCode, err := controller.Service.GetInvitingCode(ctx, phoneNumber)
	if err != nil {
		message, errRes := response.NewErrorResponse(gerror.ErrorNotFound, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(echo, message, errRes)
	}

	user := dto.User{
		Code: invitingCode,
	}

	return controller.WriteSuccess(echo, user)
}


