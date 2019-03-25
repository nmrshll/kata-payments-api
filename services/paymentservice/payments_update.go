package paymentservice

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/transact"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// UpdatePaymentIn helps unpack the params from the request body
type UpdatePaymentIn struct {
	AmountCents        *int64     `json:"amount_cents"`
	CurrencyID         *int64     `json:"currency_id"`
	BeneficiaryPartyID *int64     `json:"beneficiary_party_id"`
	DebtorPartyID      *int64     `json:"debtor_party_id"`
	PaidAt             *time.Time `json:"paid_at"`
	Reference          *string    `json:"reference"`
}

// UpdatePayment updates a Payment
func UpdatePayment(c echo.Context) error {
	return transact.Transact(func(tx *sql.Tx) error {
		ctx := context.Background()
		params := &UpdatePaymentIn{}
		if err := c.Bind(params); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		paymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
		}

		foundPayment, err := models.FindPaymentG(ctx, paymentID)
		if err != nil {
			return errors.NewNotFound(err, "failed retrieving payment")
		}

		// update each value of foundPayment if new value is not null
		if params.AmountCents != nil {
			foundPayment.AmountCents = null.Int64FromPtr(params.AmountCents)
		}
		if params.CurrencyID != nil {
			foundPayment.CurrencyID = null.Int64FromPtr(params.CurrencyID)
		}
		if params.BeneficiaryPartyID != nil {
			foundPayment.BeneficiaryPartyID = null.Int64FromPtr(params.BeneficiaryPartyID)
		}
		if params.DebtorPartyID != nil {
			foundPayment.DebtorPartyID = null.Int64FromPtr(params.DebtorPartyID)
		}
		if params.PaidAt != nil {
			foundPayment.PaymentDate = null.TimeFromPtr(params.PaidAt)
		}
		if params.Reference != nil {
			foundPayment.Reference = null.StringFromPtr(params.Reference)
		}

		_, err = foundPayment.UpdateG(ctx, boil.Infer())
		if err != nil {
			return errors.NewBadRequest(err, "failed inserting Payment")
		}

		// get the payment again to return it in the right format
		foundPayment, err = models.Payments(
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
	})
}
