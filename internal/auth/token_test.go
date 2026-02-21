package auth

import "testing"

func TestNormalizeExpiresIn(t *testing.T) {
	tests := []struct {
		input any
		want  int
	}{
		{float64(3600), 3600},
		{float64(3600000), 3600},
		{float64(0), 1},
		{float64(-1), 1},
		{"3600", 3600},
		{"3600000", 3600},
		{"bad", 3600},
		{nil, 3600},
		{true, 3600},
	}
	for _, tt := range tests {
		got := normalizeExpiresIn(tt.input)
		if got != tt.want {
			t.Errorf("normalizeExpiresIn(%v) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
