package calc

import (
	"errors"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedNum float64
		expectedErr error
	}{
		{
			name:        "easy",
			expression:  "2 + 2",
			expectedNum: 4,
			expectedErr: nil,
		},
		{
			name:        "with brackets",
			expression:  "(2+2)*3",
			expectedNum: 12,
			expectedErr: nil,
		},
		{
			name:        "invalid brackets",
			expression:  "((2+2)*3",
			expectedNum: 0,
			expectedErr: ErrInvalidBracket,
		},
		{
			name:        "invalid operands",
			expression:  "2**3",
			expectedNum: 0,
			expectedErr: ErrInvalidOperands,
		},
		{
			name:        "division by zero",
			expression:  "1/0",
			expectedNum: 0,
			expectedErr: ErrDivByZero,
		},
		{
			name:        "with double digit numbers",
			expression:  "22*3",
			expectedNum: 66,
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := Calc(test.expression)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("Name: %s\nCalc(%q): expected error %v, got %v", test.name, test.expression, test.expectedErr, err)
			}
			if val != test.expectedNum {
				t.Errorf("Name: %s\nCalc(%q): expected num %.2f, got %.2f", test.name, test.expression, test.expectedNum, val)
			}
		})
	}
}
