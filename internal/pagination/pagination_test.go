package pagination

import "testing"

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		input map[string]any
		want  bool
	}{
		{map[string]any{"has_next_page": true}, true},
		{map[string]any{"has_next_page": false}, false},
		{map[string]any{"has_next_page": "true"}, true},
		{map[string]any{"has_next_page": "false"}, false},
		{map[string]any{}, false},
		{nil, false},
	}
	for _, tt := range tests {
		got := hasNextPage(tt.input)
		if got != tt.want {
			t.Errorf("hasNextPage(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
