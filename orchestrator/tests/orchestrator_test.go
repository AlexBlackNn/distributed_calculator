package tests

import (
	"github.com/gavv/httpexpect/v2"
	"net/url"
	"orchestrator/internal/http-server/handlers/url/expression"
	"testing"
)

// TODO: get from config
const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/expression").
		WithJSON(expression.Request{
			Expression: "2+2*2",
		}).
		Expect().Status(200).JSON().Object().
		ContainsKey("id").
		ContainsKey("response")
}
