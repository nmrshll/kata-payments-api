package paymentservice

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/transact"
)

// DeletePayment deletes a Payment
func  DeletePayment(c echo.Context) error {
	return transact.Transact(func(tx *sql.Tx) error {
		ctx := context.Background()
		pathPaymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
		}

		paymentToDelete, err := models.FindPaymentG(ctx, pathPaymentID)
		if err != nil {
			return errors.NewNotFound(err, "failed finding payment to delete")
		}
		// Delete the paymentToDelete from the database
		_, err = paymentToDelete.DeleteG(ctx)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, nil)
	})
}
