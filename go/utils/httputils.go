package utils

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
)

func CreateInsecureClient(t *testing.T, baseURL string) *httpexpect.Expect {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	config := httpexpect.Config{
		BaseURL: baseURL,
		Client: &http.Client{
			Transport: transport,
			Timeout:   time.Minute * 20,
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	}
	return httpexpect.WithConfig(config)
}
