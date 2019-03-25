package api

import (
	"net/http"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nmrshll/kata-payments-api/services/paymentservice"
)

// NewServer creates a new echo server with routes, error-to-status mapping, and logging
func NewServer() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Payment routes
	e.POST("/payments", paymentservice.CreatePayment)
	e.GET("/payments", paymentservice.ListPayments)
	e.GET("/payments/:id", paymentservice.GetPayment)
	e.PATCH("/payments/:id", paymentservice.UpdatePayment)
	e.DELETE("/payments/:id", paymentservice.DeletePayment)

	return e
}

// map errors to HTTP status
func customHTTPErrorHandler(err error, c echo.Context) {
	// set default to InternalServerError
	code := http.StatusInternalServerError

	// replace statusCode if echo.HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// replace statusCode if juju/errors error type
	{
		if errors.IsNotFound(err) {
			code = http.StatusNotFound
		}
		if errors.IsBadRequest(err) {
			code = http.StatusBadRequest
		}
	}

	if err := c.JSON(code, http.StatusText(code)); err != nil {
		c.Logger().Error(err)
	}
}
