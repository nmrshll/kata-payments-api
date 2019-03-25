package paymentservice_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/convert"
	"github.com/nmrshll/kata-payments-api/libs/routes_tester"
	"github.com/nmrshll/kata-payments-api/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/nmrshll/kata-payments-api/services/paymentservice"
)

func TestPayments_Update(t *testing.T) {
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
		AmountCents:        null.Int64From(34677),
		CurrencyID:         null.Int64From(data.Currency1.ID),
		BeneficiaryPartyID: null.Int64From(data.Party1.ID),
		DebtorPartyID:      null.Int64From(data.Party2.ID),
	}

	var paymentBeforeUpdate *models.Payment

	type TestRequest struct {
		testName        string
		fnSetup         func()
		paymentToUpdate *models.Payment
		in              paymentservice.UpdatePaymentIn
		expectStatus    int
	}
	testsValidRequests := []TestRequest{
		{"valid request",
			func() {
				assert.NoError(t, payment1.InsertG(bg, boil.Infer()))
				assert.NoError(t, payment2.InsertG(bg, boil.Infer()))
			},
			payment1,
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(4500),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party1.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party2.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("another_reference"),
			},
			200},
		{"valid request 2",
			nil,
			payment1,
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(4700),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party2.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now().Add(999 * time.Hour)),
				Reference:          convert.StrToPtr("yet_another_reference"),
			},
			200},
		{"update with empty values, to check it doesn't overwrite",
			func() {
				var err error
				paymentBeforeUpdate, err = models.FindPaymentG(bg, payment1.ID)
				assert.NoError(t, err)
			},
			payment1,
			paymentservice.UpdatePaymentIn{
				Reference: convert.StrToPtr("a_different_reference"),
			},
			200},
	}
	testBadRequests := []TestRequest{
		{"update valid payment, with inexistent debtor_party_id foreign key",
			nil,
			payment2,
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(24867),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party1.ID),
				DebtorPartyID:      convert.Int64ToPtr(4),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
		{"update valid payment, with inexistent beneficiary_party_id foreign key",
			nil,
			payment1,
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(45563),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(654),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
		{"update valid payment, with inexistent currency_id foreign key",
			nil,
			payment1,
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(33251),
				CurrencyID:         convert.Int64ToPtr(984651),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party2.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
		{"try to update inexistent Payment",
			nil,
			&models.Payment{ID: 3456789},
			paymentservice.UpdatePaymentIn{
				AmountCents:        convert.Int64ToPtr(0467),
				CurrencyID:         convert.Int64ToPtr(984651),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party2.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			404},
	}

	for _, tt := range append(testsValidRequests, testBadRequests...) {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.fnSetup != nil {
				tt.fnSetup()
			}
			var out paymentservice.GetPaymentResponse

			payloadJson, err := json.Marshal(tt.in)
			assert.NoError(t, err)

			tester := routes_tester.NewTester(t, testSetup.EchoServer.Server.Handler)
			tester.AddCall(
				"updatePayment",
				"PATCH",
				fmt.Sprintf("/payments/%v", tt.paymentToUpdate.ID),
				string(payloadJson),
			).Checkers(
				routes_tester.ExpectStatus(tt.expectStatus),
				func(r *http.Response, body string, respObject interface{}) error {
					if tt.expectStatus == 200 {
						routes_tester.UnmarshalInto(t, &out)(r, body, respObject)
						// check that values were not overwritten
						if tt.in.AmountCents == nil {
							spew.Dump(out)
							assert.Equal(t, paymentBeforeUpdate.AmountCents.Int64, out.Payment.AmountCents.Int64)
						}
					}
					return nil
				},
			)
			tester.Run()
		})
	}
}
