package paymentservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nmrshll/kata-payments-api/services/paymentservice"

	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/nmrshll/kata-payments-api/libs/routes_tester"
	"github.com/nmrshll/kata-payments-api/testsetup"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

var bg = context.Background()

func TestPayments_List(t *testing.T) {
	testSetup, data := testsetup.New(t)
	defer testSetup.FnCleanup()

	// test data to insert later
	payment1 := models.Payment{
		AmountCents:        null.Int64From(98700),
		CurrencyID:         null.Int64From(data.Currency1.ID),
		BeneficiaryPartyID: null.Int64From(data.Party1.ID),
		DebtorPartyID:      null.Int64From(data.Party2.ID),
	}
	payment2 := models.Payment{
		AmountCents:        null.Int64From(98700),
		CurrencyID:         null.Int64From(data.Currency1.ID),
		BeneficiaryPartyID: null.Int64From(data.Party1.ID),
		DebtorPartyID:      null.Int64From(data.Party2.ID),
	}

	type TestRequest struct {
		testName            string
		fnSetup             func()
		expectNumberResults int
		expectStatus        int
	}
	testsValidRequests := []TestRequest{
		{"valid request, only 1 item present",
			func() {
				err := payment1.InsertG(bg, boil.Infer())
				assert.NoError(t, err)
			},
			1,
			200},
		{"valid request, now 2 items present",
			func() {
				err := payment2.InsertG(bg, boil.Infer())
				assert.NoError(t, err)
			},
			2,
			200},
	}
	testBadRequests := []TestRequest{}

	for _, tt := range append(testsValidRequests, testBadRequests...) {
		t.Run(tt.testName, func(t *testing.T) {
			tt.fnSetup()
			var out paymentservice.ListPaymentsResponse

			tester := routes_tester.NewTester(t, testSetup.EchoServer.Server.Handler)
			tester.AddCall(
				"listPayments",
				"GET",
				fmt.Sprintf("/payments"),
				"",
			).Checkers(
				routes_tester.ExpectStatus(tt.expectStatus),
				routes_tester.ExpectJSONFields("payments"),
				routes_tester.UnmarshalInto(t, &out,
					func() { assert.Len(t, out.Payments, tt.expectNumberResults) },
				),
			)
			tester.Run()
		})
	}
}
