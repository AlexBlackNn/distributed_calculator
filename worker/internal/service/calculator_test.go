package service

import (
	"fmt"
	"testing"
	transport "worker/internal/transport"
)

func TestCalculatorOperation(t *testing.T) {
	var tests = []struct {
		expression string
		result     int
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

	for _, test := range tests {
		name := fmt.Sprintf("case(%v,%v)", test.expression, test.result)
		calculator := New()
		t.Run(name, func(t *testing.T) {
			requestMessage := transport.RequestMessage{Id: "1231", Operation: test.expression}
			result, err := calculator.Start(requestMessage)
			if err != nil {
				t.Errorf("got %v, want %v", result, test.result)
			}
			if result != test.result {
				t.Errorf("got %v, want %v", result, test.result)
			}
		})
	}
}
