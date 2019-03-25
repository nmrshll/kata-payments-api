package paymentservice_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nmrshll/kata-payments-api/libs/convert"
	"github.com/nmrshll/kata-payments-api/libs/routes_tester"
	"github.com/nmrshll/kata-payments-api/testsetup"
	"github.com/stretchr/testify/assert"

	"github.com/nmrshll/kata-payments-api/services/paymentservice"
)

func TestPayments_Create(t *testing.T) {
	testSetup, data := testsetup.New(t)
	defer testSetup.FnCleanup()

	type TestRequest struct {
		testName     string
		fnSetup      func()
		in           paymentservice.CreatePaymentIn
		expectStatus int
	}
	testsValidRequests := []TestRequest{
		{"valid request",
			nil,
			paymentservice.CreatePaymentIn{
				AmountCents:        convert.Int64ToPtr(98700),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party1.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party2.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			200},
	}
	testBadRequests := []TestRequest{
		{"inexistent debtor_party_id foreign key",
			nil,
			paymentservice.CreatePaymentIn{
				AmountCents:        convert.Int64ToPtr(98700),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party1.ID),
				DebtorPartyID:      convert.Int64ToPtr(4),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
		{"inexistent beneficiary_party_id foreign key",
			nil,
			paymentservice.CreatePaymentIn{
				AmountCents:        convert.Int64ToPtr(98700),
				CurrencyID:         convert.Int64ToPtr(data.Currency1.ID),
				BeneficiaryPartyID: convert.Int64ToPtr(654),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
		{"inexistent currency_id foreign key",
			nil,
			paymentservice.CreatePaymentIn{
				AmountCents:        convert.Int64ToPtr(98700),
				CurrencyID:         convert.Int64ToPtr(984651),
				BeneficiaryPartyID: convert.Int64ToPtr(data.Party2.ID),
				DebtorPartyID:      convert.Int64ToPtr(data.Party1.ID),
				PaidAt:             convert.TimeToPtr(time.Now()),
				Reference:          convert.StrToPtr("a_reference"),
			},
			400},
	}

	for _, tt := range append(testsValidRequests, testBadRequests...) {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.fnSetup != nil {
				tt.fnSetup()
			}

			payloadJson, err := json.Marshal(tt.in)
			assert.NoError(t, err)

			tester := routes_tester.NewTester(t, testSetup.EchoServer.Server.Handler)
			tester.AddCall(
				"createPayment",
				"POST",
				fmt.Sprintf("/payments"),
				string(payloadJson),
			).Checkers(
				routes_tester.ExpectStatus(tt.expectStatus),
			)
			tester.Run()
		})
	}
}
