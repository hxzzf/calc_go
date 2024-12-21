package calculation_test

import (
	"math"
	"testing"

	"github.com/hxzzf/calc_go/pkg/calculation"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
		errorMsg   string
	}{
		{
			name:       "simple addition",
			expression: "2 + 3",
			want:       5,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "simple multiplication",
			expression: "4 * 5",
			want:       20,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "parentheses",
			expression: "(2 + 3) * 4",
			want:       20,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "complex expression",
			expression: "2.5 * (3 + 4) / 2",
			want:       8.75,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "division by zero",
			expression: "1 / 0",
			want:       0,
			wantErr:    true,
			errorMsg:   "division by zero",
		},
		{
			name:       "invalid token",
			expression: "2 + a",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid token",
		},
		{
			name:       "mismatched parentheses",
			expression: "(2 + 3",
			want:       0,
			wantErr:    true,
			errorMsg:   "mismatched parentheses",
		},
		{
			name:       "decimal numbers",
			expression: "1.5 + 2.5",
			want:       4,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "operator precedence",
			expression: "2 + 3 * 4",
			want:       14,
			wantErr:    false,
			errorMsg:   "",
		},
		{
			name:       "empty expression",
			expression: "",
			want:       0,
			wantErr:    true,
			errorMsg:   "empty expression",
		},
		{
			name:       "only spaces",
			expression: "    ",
			want:       0,
			wantErr:    true,
			errorMsg:   "empty expression",
		},
		{
			name:       "consecutive operators",
			expression: "2 ++ 3",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid expression: consecutive operators",
		},
		{
			name:       "missing operator",
			expression: "2 3",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid expression",
		},
		{
			name:       "starts with operator",
			expression: "+ 2 + 3",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid expression",
		},
		{
			name:       "ends with operator",
			expression: "2 + 3 +",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid expression",
		},
		{
			name:       "too many closing parentheses",
			expression: "(2 + 3))",
			want:       0,
			wantErr:    true,
			errorMsg:   "mismatched parentheses",
		},
		{
			name:       "invalid decimal number",
			expression: "2.5.6 + 1",
			want:       0,
			wantErr:    true,
			errorMsg:   "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculation.Calc(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(got, tt.want) {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errorMsg {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.errorMsg)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-10
}
