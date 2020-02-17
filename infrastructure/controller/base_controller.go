package controller

import (
	"g-tech.com/infrastructure/response"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type BaseController struct {

}

/**
 * Returns a success response
 */
func (controller *BaseController) WriteSuccess(c echo.Context, v interface{}) error {
	response := response.Response{
		Message: "Success",
		Data:    v,
	}

	// Log response
	//logger.Info(util.ToJSON(response))

	// Return
	return c.JSON(http.StatusOK, response)
}

/**
 * Returns a success response without content
 */
func (controller *BaseController) WriteSuccessEmptyContent(c echo.Context) error {
	response := response.Response{
		Message: "Success",
		Data:    nil,
	}

	// Log response
	//logger.Info(util.ToJSON(response))

	// Return
	return c.JSON(http.StatusOK, response)
}

/**
 * Returns an error as bad request (client-side error)
 */
func (controller *BaseController) WriteBadRequest(c echo.Context, message string, errorRes response.ErrorResponse)  error {
	return controller.writeError(c, http.StatusBadRequest, message, errorRes)
}

/**
 * Return an error as NotFound (client-side error
 */
func (c *BaseController) WriteStatusNotFound(e echo.Context, message string, errorRes response.ErrorResponse)  error {
	return c.writeError(e, http.StatusNotFound, message, errorRes)
}

/**
 * Return an error as Conflict (client-side error
 */
func (c *BaseController) WriteStatusConflict(e echo.Context, message string, errorRes response.ErrorResponse)  error {
	return c.writeError(e, http.StatusConflict, message, errorRes)
}

/**
 * Redirect an error as internal server error (server-side error)
 */
func (controller *BaseController) WriteInternalServerError(c echo.Context, message string, errorRes response.ErrorResponse)  error {
	return controller.writeError(c, http.StatusInternalServerError, message, errorRes)
}

/**
 * Returns an error response
 */
func (controller *BaseController) writeError(c echo.Context, statusCode int, message string, err response.ErrorResponse)  error {
	response := response.Response{
		Message: message,
		Data:    err,
	}

	// Log error
	//logger.Error(util.ToJSON(response))

	// Return
	return c.JSON(statusCode, response)
}

/**
 * Validates model before do something
 */
func (controller *BaseController) IsValid(m interface{}) (bool, error) {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

