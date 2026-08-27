package keyring

import "testing"

func TestIsReference(t *testing.T) {
	tests := []struct {
		raw  string
		want bool
	}{
		{"keyring://", true},
		{"keyring://foo", false},
		{"", false},
		{"https://example.com", false},
	}
	for _, tc := range tests {
		got := IsReference(tc.raw)
		if got != tc.want {
			t.Errorf("IsReference(%q) = %v, want %v", tc.raw, got, tc.want)
		}
	}
}
