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

const (
	host = "localhost:8080"
)

func TestExpression_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/expression").
		WithJSON(expression.Request{
			Expression: "2+2*1",
		}).
		Expect().Status(http.StatusOK).JSON().Object().
		ContainsKey("id").
		ContainsKey("response")
}

func TestWholePathOfCalculation2_HappyPath(t *testing.T) {
	var tests = []struct {
		expression string
		result     float64
	}{
		{"1+2", 3},
		{"1+2*2", 5},
		{"(1+2)*2", 6},
		{"(1+2)*2/3", 2},
		{"(1+2)*2/3*(1-2)", -2},
		{"2+(1+2)*2/3*(1-2)", 0},
		{"2*2/2", 2},
		{"1", 1},
		{"-1", -1},
		{"0-2", -2},
		{"0+2", 2},
	}

	interactor := Interactor{
		client: &http.Client{Timeout: 3 * time.Second},
		url:    "http://localhost:8080/expression",
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v,%v)", test.expression, test.result)
		t.Run(name, func(t *testing.T) {
			request := expression.Request{Expression: test.expression}
			response := expression.Response{}
			err := interactor.PostExpression(request, &response)
			if err != nil {
				t.Errorf(err.Error())
			}
			time.Sleep(4 * time.Second)
			err = interactor.GetResult(&response)
			if err != nil {
				t.Errorf(err.Error())
			}
			if response.Response.Result != test.result {
				t.Errorf("got %v, want %v", response.Response.Result, test.result)
			}
		})
	}
}
