package magcal

import (
	"testing"
)

func TestDefaultState(t *testing.T) {
	s := DefaultState()
	v := vector{1, 2, 3}
	w := s.apply(v)
	if !v.close(w, 3) {
		t.Fatalf(`want %v got %v`, v, w)
	}
}

func TestCustomState(t *testing.T) {
	cal := []float32{
		0.5, 0.0, 0.0, // offset
		1.0, 0.0, 0.0, // covariance matrix
		0.0, 1.2, 0.0,
		0.0, 0.0, 1.0,
	}
	s := NewState(cal)
	v := vector{1, 2, 3}
	expect := vector{0.5, 2.4, 3}

	actual := s.apply(v)
	if !actual.close(expect, 3) {
		t.Fatalf(`want %v got %v`, expect, actual)
	}
}
