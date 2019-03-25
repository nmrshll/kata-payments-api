package paymentservice

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/transact"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// CreatePaymentIn helps unpack the parameters from the request body
type CreatePaymentIn struct {
	AmountCents        *int64     `json:"amount_cents"`
	CurrencyID         *int64     `json:"currency_id"`
	BeneficiaryPartyID *int64     `json:"beneficiary_party_id"`
	DebtorPartyID      *int64     `json:"debtor_party_id"`
	PaidAt             *time.Time `json:"paid_at"`
	Reference          *string    `json:"reference"`
}

// CreatePayment creates a Payment
func CreatePayment(c echo.Context) error {
	return transact.Transact(func(tx *sql.Tx) error {
		ctx := context.Background()
		params := &CreatePaymentIn{}
		if err := c.Bind(params); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		newPayment := models.Payment{
			AmountCents:        null.Int64FromPtr(params.AmountCents),
			CurrencyID:         null.Int64FromPtr(params.CurrencyID),
			BeneficiaryPartyID: null.Int64FromPtr(params.BeneficiaryPartyID),
			DebtorPartyID:      null.Int64FromPtr(params.DebtorPartyID),
			PaymentDate:        null.TimeFromPtr(params.PaidAt),
			Reference:          null.StringFromPtr(params.Reference),
		}
		err := newPayment.InsertG(ctx, boil.Infer())
		if err != nil {
			return errors.NewBadRequest(err, "failed inserting Payment")
		}

		return c.JSON(http.StatusOK, newPayment)
	})
}
