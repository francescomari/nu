package paths

import "testing"

func TestComponents(t *testing.T) {
	tests := []struct {
		path       string
		components []string
	}{
		{"/", nil},
		{"//", nil},
		{"/a", []string{"a"}},
		{"//a", []string{"a"}},
		{"//a/", []string{"a"}},
		{"/aa", []string{"aa"}},
		{"/aa/bb", []string{"aa", "bb"}},
	}

	for _, tt := range tests {
		if components, err := Components(tt.path); err != nil {
			t.Errorf("path %v: unexpected error: %v\n", tt.path, err)
		} else if !areStringsEqual(components, tt.components) {
			t.Errorf("path %v: expected %v, got %v\n", tt.path, tt.components, components)
		}
	}
}

func areStringsEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
