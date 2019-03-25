package contract_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"stash.ovh.net/agora/contract/testsutils"
	"stash.ovh.net/agora/contract/utils/random"
	models "stash.ovh.net/agora/go-db-models"
	"stash.ovh.net/agora/go-db-models/helper"
	"stash.ovh.net/golang/httptester"
)

func TestGetContract(t *testing.T) {
	fnCleanup, dbp, ginRouter, _ := testsutils.Init(t)
	defer fnCleanup()

	contractURLGen := random.NewURLGenerator().OnGenerate(func(oldURL, newURL string) {
		httpmock.Activate()
		httpmock.RegisterResponder("GET", newURL, func(r *http.Request) (*http.Response, error) {
			contractResponse := http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("Hello-world contract " + random.RandStringBytes(8))),
			}
			return &contractResponse, nil
		})
	})

	contractInsertedV1, err := helper.InsertContract(dbp, "my_contract_name", *contractURLGen.Next())
	assert.NoError(t, err)
	assert.NotNil(t, contractInsertedV1)
	contractInsertedV2, err := helper.InsertContract(dbp, "my_contract_name", *contractURLGen.Next())
	assert.NoError(t, err)
	assert.NotNil(t, contractInsertedV2)

	tester := httptester.NewTester(t, ginRouter)
	tester.AddCall(
		"getContract",
		"GET",
		fmt.Sprintf("/contract/name/%s", *contractInsertedV1.Name),
		"",
	).ResponseObject(&models.Contract{}).Checkers(
		httptester.ExpectStatus(200),
		verifyContract(t, contractInsertedV2.ID),
	)
	tester.Run()

	tester = httptester.NewTester(t, ginRouter)
	tester.AddCall(
		"getContract",
		"GET",
		fmt.Sprintf("/contract/name/%s", "invalid_name"),
		"",
	).ResponseObject(&models.Contract{}).Checkers(
		httptester.ExpectStatus(404),
	)
	tester.Run()
}

func verifyContract(t *testing.T, contractId int64) httptester.Checker {
	return func(r *http.Response, body string, respObject interface{}) error {
		ContractCreated := respObject.(*models.Contract)
		assert.Equal(t, contractId, ContractCreated.ID, "ContractID is different")
		return nil
	}
}
