package testground

import "testing"

func AssertStringNotChanged(t *testing.T, getter func(t *testing.T) string) func() {
	t.Helper()
	before := getter(t)

	return func() {
		t.Helper()
		after := getter(t)
		if before != after {
			t.Fatalf("expected %q, but got %q", before, after)
		}
	}
}
