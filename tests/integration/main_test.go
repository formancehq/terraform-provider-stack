package integration_test

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/otlp"
)

var transport http.RoundTripper

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		transport = httpclient.NewDebugHTTPTransport(
			otlp.NewRoundTripper(http.DefaultTransport, true),
		)
	} else {
		transport = otlp.NewRoundTripper(http.DefaultTransport, false)
	}
	code := m.Run()

	os.Exit(code)
}
