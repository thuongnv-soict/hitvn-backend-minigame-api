package healthcheck

import (
	"database/sql"
	"g-tech.com/module/healthcheck/controller"
	"g-tech.com/module/healthcheck/respository"
	"g-tech.com/module/healthcheck/service"
	"github.com/labstack/echo"
	"time"
)

var mHealthCheckController *controller.HealthCheckController

/**
 * Initializes module
 */
func Initialize(e *echo.Echo, dbContext *sql.DB, timeout time.Duration) {
	healthCheckRepository := respository.NewHealthCheckRepository(dbContext, timeout)
	healthCheckService :=  service.NewHealthCheckService(healthCheckRepository)

	mHealthCheckController = controller.NewHealthCheckController(healthCheckService)

	// New router
	initRouter(e)
}

/**
 * Initializes router
 */
func initRouter(e *echo.Echo) {
	e.GET("/game/api/v1.0/status", mHealthCheckController.GetStatus)
	e.GET("/game/api/v1.0/database/status", mHealthCheckController.GetSQLStatus)
}


