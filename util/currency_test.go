package util

import "testing"

func TestIsSupportedCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		currency string
		want     bool
	}{
		{"USD", USD, true},
		{"EUR", EUR, true},
		{"CAD", CAD, true},
		{"Unsupported", "GBP", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedCurrency(tt.currency); got != tt.want {
				t.Errorf("IsSupportedCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}
