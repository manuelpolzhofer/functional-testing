package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
)

func createDocumentForNFT(t *testing.T) *httpexpect.Object {
	e := utils.GetInsecureClient(t, utils.NODE1)

	// createDocument
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"invoice_number": "12324",
			"due_date":       "2018-09-26T23:12:37.902198664Z",
			"gross_amount":   "40",
			"currency":       "USD",
			"net_amount":     "40",
			"document_type":  "invoice",
		},
		"collaborators": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj := CreateDocument(t, utils.INVOICE, e, payload)

	return obj

}

func TestPaymentObligationMint_successful(t *testing.T) {

	const tokenIdLength = 77
	utils.GetInsecureClient(t, utils.NODE1)

	expectedNode1 := utils.GetInsecureClient(t, utils.NODE1)

	docObj := createDocumentForNFT(t)
	documentId := docObj.Value("header").Object().Value("document_id").String().Raw()
	test := struct {
		errorMsg   string
		httpStatus int
		payload    map[string]interface{}
	}{
		"",
		http.StatusOK,
		map[string]interface{}{

			"identifier":      documentId,
			"registryAddress": utils.GetPaymentObigationAddress(),
			"depositAddress":  "0xf72855759a39fb75fc7341139f5d7a3974d4da08", // dummy address
			"proofFields":     []string{"invoice.gross_amount", "invoice.currency", "invoice.due_date", "collaborators[0]"},
		},
	}

	response := PostTokenMint(expectedNode1, test.httpStatus, test.payload)
	assert.True(t, len(response.Value("token_id").String().Raw()) >= tokenIdLength, "successful tokenId should have length 77")

}

func TestPaymentObligationMint_errors(t *testing.T) {
	expectedNode1 := utils.GetInsecureClient(t, utils.NODE1)

	tests := []struct {
		errorMsg   string
		httpStatus int
		payload    map[string]interface{}
	}{
		{
			"RegistryAddress is not a valid Ethereum address",
			http.StatusInternalServerError,
			map[string]interface{}{

				"registryAddress": "0x123",
			},
		},
		{
			"DepositAddress is not a valid Ethereum address",
			http.StatusInternalServerError,
			map[string]interface{}{

				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08", //dummy address
				"depositAddress":  "abc",
			},
		},
		{
			"no service exists for provided documentID",
			http.StatusInternalServerError,
			map[string]interface{}{

				"identifier":      "0x12121212",
				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08", //dummy address
				"depositAddress":  "0xf72855759a39fb75fc7341139f5d7a3974d4da08", //dummy address
			},
		},
	}

	for _, test := range tests {
		httpObj := PostTokenMint(expectedNode1, test.httpStatus, test.payload)
		httpObj.Value("message").String().Contains(test.errorMsg)

	}

}

func PostTokenMint(e *httpexpect.Expect, httpStatus int, payload map[string]interface{}) *httpexpect.Object {
	resp := e.POST("/token/mint").
		WithHeader("accept", "application/json").
		WithHeader("Content-Type", "application/json").
		WithJSON(payload).
		Expect().Status(httpStatus)

	httpObj := resp.JSON().Object()
	return httpObj
}
