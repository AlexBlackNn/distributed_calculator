package tests

import (
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"orchestrator/internal/http-server/handlers/url/expression"
	"testing"
	"time"
)

// TODO: get from config
const (
	host = "localhost:8080"
)

func TestWholePathOfCalculation_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	obj := e.POST("/expression").
		WithJSON(expression.Request{
			Expression: "2+2*1",
		}).
		Expect().Status(http.StatusOK).JSON().Object().
		ContainsKey("id").
		ContainsKey("response")

	var uidReturned string
	obj.Value("id").Decode(&uidReturned)
	fmt.Print(uidReturned)
	time.Sleep(1 * time.Second)

	e.GET("/expression/" + uidReturned).Expect().
		Status(http.StatusOK).JSON().Object()

}
