package tests

import (
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
)

func TestPing(t *testing.T) {
	e := utils.GetInsecureClient(t, utils.NODE1)
	obj := e.GET("/ping").
		Expect().
		Status(http.StatusOK)
	assertOkResponse(t, obj)
	obj.JSON().Object().Value("network").Equal(utils.Network)
	obj.JSON().Object().ContainsKey("version")
}
