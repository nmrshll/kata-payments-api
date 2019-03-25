package paymentservice_test

import (
	"fmt"
	"net/http"
	"testing"

	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/routes_tester"
	"github.com/nmrshll/kata-payments-api/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestPayments_Delete(t *testing.T) {
	testSetup, data := testsetup.New(t)
	defer testSetup.FnCleanup()

	// test data to insert later
	payment1 := &models.Payment{
		AmountCents:        null.Int64From(98700),
		CurrencyID:         null.Int64From(data.Currency1.ID),
		BeneficiaryPartyID: null.Int64From(data.Party1.ID),
		DebtorPartyID:      null.Int64From(data.Party2.ID),
	}
	payment2 := &models.Payment{
		AmountCents:        null.Int64From(98700),
		CurrencyID:         null.Int64From(data.Currency1.ID),
		BeneficiaryPartyID: null.Int64From(data.Party1.ID),
		DebtorPartyID:      null.Int64From(data.Party2.ID),
	}

	type TestRequest struct {
		testName        string
		fnSetup         func()
		paymentIn       *models.Payment
		expectRemaining int
		expectStatus    int
	}
	testsValidRequests := []TestRequest{
		{"valid request, 2 items present, delete the first one",
			func() {
				assert.NoError(t, payment1.InsertG(bg, boil.Infer()))
				assert.NoError(t, payment2.InsertG(bg, boil.Infer()))
			},
			payment1, 1, 200},
		{"valid request, 1 item left, delete it",
			nil,
			payment2, 0, 200},
	}
	testBadRequests := []TestRequest{
		{"try to delete payment1 again",
			nil,
			payment1, 0, 404},
		{"try to delete payment2 again",
			nil,
			payment1, 0, 404},
		{"try to delete an inexistent payment",
			nil,
			&models.Payment{ID: 456789}, 0, 404},
	}

	for _, tt := range append(testsValidRequests, testBadRequests...) {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.fnSetup != nil {
				tt.fnSetup()
			}

			tester := routes_tester.NewTester(t, testSetup.EchoServer.Server.Handler)
			tester.AddCall(
				"deletePayment",
				"DELETE",
				fmt.Sprintf("/payments/%v", tt.paymentIn.ID),
				"",
			).Checkers(
				routes_tester.ExpectStatus(tt.expectStatus),
				func(r *http.Response, body string, respObject interface{}) error {
					if tt.expectStatus == 200 {
						foundPayments, err := models.Payments().AllG(bg)
						assert.NoError(t, err)
						assert.Len(t, foundPayments, tt.expectRemaining)
					}
					return nil
				},
			)
			tester.Run()
		})
	}
}
