package paymentservice_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nmrshll/kata-payments-api/services/paymentservice"

	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/routes_tester"
	"github.com/nmrshll/kata-payments-api/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestPayments_Get(t *testing.T) {
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
		testName     string
		fnSetup      func()
		paymentIn    *models.Payment
		expectStatus int
	}
	testsValidRequests := []TestRequest{
		{"valid request, item 1",
			func() {
				err := payment1.InsertG(bg, boil.Infer())
				assert.NoError(t, err)
			},
			payment1,
			200},
		{"valid request, item 2",
			func() {
				err := payment2.InsertG(bg, boil.Infer())
				assert.NoError(t, err)
			},
			payment2,
			200},
	}
	testBadRequests := []TestRequest{
		{"inexistent item",
			nil,
			&models.Payment{
				ID: 987654,
			},
			404},
	}

	for _, tt := range append(testsValidRequests, testBadRequests...) {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.fnSetup != nil {
				tt.fnSetup()
			}
			var out paymentservice.GetPaymentResponse

			tester := routes_tester.NewTester(t, testSetup.EchoServer.Server.Handler)
			tester.AddCall(
				"getPayment",
				"GET",
				fmt.Sprintf("/payments/%v", tt.paymentIn.ID),
				"",
			).Checkers(
				routes_tester.ExpectStatus(tt.expectStatus),
				func(r *http.Response, body string, respObject interface{}) error {
					if tt.expectStatus == 200 {
						assert.NoError(t, routes_tester.ExpectJSONFields("payment")(r, body, respObject))
						assert.NoError(t, routes_tester.UnmarshalInto(t, &out)(r, body, respObject))

						assert.NotZero(t, out.Payment.ID)
						assert.Equal(t, tt.paymentIn.ID, out.Payment.ID)
						assert.Equal(t, tt.paymentIn.BeneficiaryPartyID, out.Payment.BeneficiaryPartyID)
						assert.Equal(t, tt.paymentIn.DebtorPartyID, out.Payment.DebtorPartyID)
						assert.Equal(t, tt.paymentIn.AmountCents, out.Payment.AmountCents)
					}
					return nil
				},
			)
			tester.Run()
		})
	}
}
