package paymentservice

import (
	"context"
	"net/http"
	"strconv"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	models "github.com/nmrshll/kata-payments-api/generated-models"
)

// PaymentResponse is a sub-object of GetPaymentResponse
// it helps expand sub-fields (relations) from only foreign keys to the full sub-object
type PaymentResponse struct {
	*models.Payment  `json:"payment"`
	Currency         *models.Currency `json:"currency"`
	BeneficiaryParty *models.Party    `json:"beneficiary_party"`
	DebtorParty      *models.Party    `json:"debtor_party"`
}

// GetPaymentResponse is the shape of the API response for GetPayment
type GetPaymentResponse struct {
	Payment PaymentResponse `json:"payment"`
}

// GetPayment returns details for one Payment
func  GetPayment(c echo.Context) error {
	ctx := context.Background()
	paymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	foundPayment, err := models.Payments(
		models.PaymentWhere.ID.EQ(paymentID),
		qm.Load(models.PaymentRels.Currency),
		qm.Load(models.PaymentRels.BeneficiaryParty),
		qm.Load(models.PaymentRels.DebtorParty),
	).OneG(ctx)
	if err != nil {
		return errors.NewNotFound(err, "failed retrieving payment")
	}

	resp := GetPaymentResponse{
		Payment: PaymentResponse{
			Payment:          foundPayment,
			Currency:         foundPayment.R.Currency,
			BeneficiaryParty: foundPayment.R.BeneficiaryParty,
			DebtorParty:      foundPayment.R.DebtorParty,
		}}

	return c.JSON(http.StatusOK, resp)
}
