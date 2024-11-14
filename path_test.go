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
		"Does not start with separator": {
			in:  "A/B",
			out: []string{"A", "B"},
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
		"Not a filesystem path, dots do not have special treatment": {
			in:  "../../A/B/",
			out: []string{"..", "..", "A", "B"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			o := slicePath(test.in)
			if !slices.Equal(o, test.out) {
				t.Errorf("unexpected output in %s:\n[GOT]\n%+v\n\n[EXPECTED]\n%+v", name, o, test.out)
			}
		})
	}
}
