package paymentservice

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// ListPaymentsResponse is the struct for holding and serializing all the fields we want in the ListPayments API response
type ListPaymentsResponse struct {
	Payments []PaymentResponse `json:"payments"`
}

// ListPayments lists all payments
func  ListPayments(c echo.Context) error {
	ctx := context.Background()

	foundPayments, err := models.Payments(
		qm.Load(models.PaymentRels.Currency),
		qm.Load(models.PaymentRels.BeneficiaryParty),
		qm.Load(models.PaymentRels.DebtorParty),
	).All(ctx, boil.GetContextDB())
	if err != nil {
		return err
	}

	resp := ListPaymentsResponse{
		// init slice to marshal to empty array, not to null
		Payments: make([]PaymentResponse, 0),
	}
	for _, p := range foundPayments {
		resp.Payments = append(resp.Payments, PaymentResponse{
			Payment:          p,
			Currency:         p.R.Currency,
			BeneficiaryParty: p.R.BeneficiaryParty,
			DebtorParty:      p.R.DebtorParty,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
