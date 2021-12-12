package metrics

import (
	"math"
	"testing"
)

func TestF1Score(t *testing.T) {
	p := 0.8
	r := 0.7
	beta := 1.0
	want := 0.74666666666
	got := F1Score(p, r, beta)

	if math.Abs(got-want) > 0.0001 {
		t.Errorf("got %.2f, want %.2f", got, want)
	}
}
