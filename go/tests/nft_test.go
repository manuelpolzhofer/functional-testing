package tests

import (
	"fmt"
	"net/http"
	"testing"

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
		},
		"collaborators": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj := CreateDocument(t, utils.INVOICE, e, payload)

	return obj

}

func TestPaymentObligationMint_Errors(t *testing.T) {
	expectedNode1 := utils.GetInsecureClient(t, utils.NODE1)

	//docObj := createDocumentForNFT(t)
	//documentId := docObj.Value("header").Object().Value("document_id").String().Raw()

	documentId := "0xc7f569f394c9da863949b26891db5f781f7d0f8cbb60a3748518a1f8c1803117"
	fmt.Println("document Id")
	fmt.Println(documentId)


	inv := GetDocument(t,utils.INVOICE,expectedNode1,documentId)
	fmt.Println(inv.Raw())

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

				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
				"depositAddress":  "abc",
			},
		},
		{
			"no service exists for provided documentID",
			http.StatusInternalServerError,
			map[string]interface{}{

				"identifier":      "0x12121212",
				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
				"depositAddress":  "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
			},
		},
		{
			"proof_fields should contain a collaborator",
			http.StatusInternalServerError,
			map[string]interface{}{

				"identifier": documentId,
				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
				"depositAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",

			},

		},
		{
			"proof_fields should contain a collaborator",
			http.StatusInternalServerError,
			map[string]interface{}{

				"identifier": documentId,
				"registryAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
				"depositAddress": "0xf72855759a39fb75fc7341139f5d7a3974d4da08",
				"proofFields":    []string{"gross_amount", "currency", "due_date", "document_type", "collaborators[0]"},

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
