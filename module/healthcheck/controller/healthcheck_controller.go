package controller

import (
	"context"
	"g-tech.com/gerror"
	"g-tech.com/infrastructure/controller"
	"g-tech.com/infrastructure/response"
	"g-tech.com/infrastructure/util"
	"g-tech.com/module/healthcheck/service"
	"github.com/labstack/echo"
)

type HealthCheckController struct {
	controller.BaseController
	Service service.IHealthCheckService
}

/**
 * Returns a new HealthCheckController
 */
func NewHealthCheckController(service service.IHealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		Service: service,
	}
}

/**
 * Returns status
 */
func (controller *HealthCheckController) GetStatus(c echo.Context) error {
	return controller.WriteSuccessEmptyContent(c)
}

/**
 * Returns Couchbase status
 */
func (controller *HealthCheckController) GetSQLStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Test connect to Couchbase
	err := controller.Service.GetSQLStatus(ctx)
	if err != nil {
		msg, errResponse := response.NewErrorResponse(gerror.ErrorConnect, err.Error(), util.FuncName())
		return controller.WriteInternalServerError(c, msg, errResponse)
	}

	// Return
	return controller.WriteSuccessEmptyContent(c)
}