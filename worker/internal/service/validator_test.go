package service

import (
	"fmt"
	"testing"
)

func TestValidOperation(t *testing.T) {
	var tests = []struct {
		expression string
		result     bool
	}{
		{"1+2", true},
		{"1+2*2", true},
		{"(1+2)*2", true},
		{"(1+2)*2/3", true},
		{"(1+2)*2/3*(1-2)", true},
		{"2+(1+2)*2/3*(1-2)", true},
		{"2*2/2", true},
		{"1", true},
		{"-1", true},
		{"0-2", true},
		{"0+2", true},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v,%v)", test.expression, test.result)
		t.Run(name, func(t *testing.T) {

			result := VerifyExpression(test.expression)

			if result != test.result {
				t.Errorf("got %v, want %v", result, test.result)
			}
		})
	}
}

func TestBadOperation(t *testing.T) {
	var tests = []struct {
		expression string
		result     bool
	}{
		{"1+2)", false},
		{"(1+2*2", false},
		{"1+2)*2", false},
		{"1++2*2/3", false},
		{"(1+2-*2/3*1-2)", false},
		{"2+1+2-*2/3*(1-2)", false},
		{"2*2/(2", false},
		{"1)", false},
		{"*-1", false},
		{"0)-2", false},
		{"0+(2", false},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v,%v)", test.expression, test.result)
		t.Run(name, func(t *testing.T) {

			result := VerifyExpression(test.expression)

			if result != test.result {
				t.Errorf("got %v, want %v", result, test.result)
			}
		})
	}
}
