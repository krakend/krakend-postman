package postman

import (
	"slices"
	"testing"
)

func TestSlicePath(t *testing.T) {
	tests := map[string]struct {
		in  string
		out []string
	}{
		"No elements": {
			in:  "/",
			out: []string{},
		},
		"Single element": {
			in:  "/A",
			out: []string{"A"},
		},
		"Multiple elements": {
			in:  "/A/B",
			out: []string{"A", "B"},
		},
		"Trailing separator": {
			in:  "/A/B/",
			out: []string{"A", "B"},
		},
		"Double separator": {
			in:  "/A//B",
			out: []string{"A", "B"},
		},
		"Non alphanumeric chars": {
			in:  "/A1/B with spaces/.hidden/_/-&",
			out: []string{"A1", "B with spaces", ".hidden", "_", "-&"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			o := SlicePath(test.in)
			if !slices.Equal(o, test.out) {
				t.Errorf("unexpected output in %s:\n[GOT]\n%+v\n\n[EXPECTED]\n%+v", name, o, test.out)
			}
		})
	}
}
