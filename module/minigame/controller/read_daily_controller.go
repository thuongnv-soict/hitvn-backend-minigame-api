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

type ReadDailyController struct {
	controller.BaseController
	Service     service.IReadDailyService
}

func NewReadDailyController(ReadDailyService service.IReadDailyService) *ReadDailyController{
	return &ReadDailyController{
		Service: ReadDailyService,
	}
}


func (controller *ReadDailyController) CreateNewReadDaily(echo echo.Context) error {
	// 1. Log ip
	logger.Trace("From %s call to %s", echo.RealIP(), util.FuncName())

	// 1. get param
	user := dto.User{}
	err := echo.Bind(&user)

	// 2. Defines context
	ctx := echo.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// 3. Retrieve data
	errorCode, err := controller.Service.CreateNewReadDaily(ctx, user)
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
