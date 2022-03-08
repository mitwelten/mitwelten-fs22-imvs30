package aggregator

import "testing"

func TestAdd(t *testing.T) {

	got := Min(4, 6)
	want := 4

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
