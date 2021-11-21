package magcal

import (
	"math/rand"
	"testing"
)

func TestMulSimple(t *testing.T) {
	m := matrix(make([]float32, 9))
	v := make([]float32, 3)
	for i := 0; i < 3; i++ {
		v[i] = 3*rand.Float32() - 1.5
		m[i*3+i] = 1
	}
	w := m.mul(v)
	if !w.close(v, 3) {
		t.Errorf(`want %v got %v`, v, w)
	}
}

func TestTrans(t *testing.T) {
	m := matrix(make([]float32, 9))
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i*3+j] = 3*rand.Float32() - 1.5
		}
	}
	mm := m.trans()
	mmm := mm.trans()
	if !mmm.close(m, 3) {
		print(m.string())
		print(mmm.string())
		t.Errorf(`trans not equal`)
	}
}

func TestInvSimple(t *testing.T) {
	m := matrix(make([]float32, 9))
	a := vector(make([]float32, 3))
	for i := 0; i < 3; i++ {
		a[i] = 3*rand.Float32() - 1.5
		m[i*3+i] = 1 + float32(i)
	}
	b := m.mul(a)
	w := m.inv()
	actual := w.mul(b)
	expected := a
	if !actual.close(expected, 3) {
		print(m.string())
		println(a.string())
		println()

		println(b.string())
		println()

		println(w.string())
		println()

		t.Errorf(`want %v got %v`, expected, actual)
	}
}

func TestInv(t *testing.T) {
	m := matrix(make([]float32, 9))
	v := vector(make([]float32, 3))
	for i := 0; i < 3; i++ {
		v[i] = 3*rand.Float32() - 1.5
		for j := 0; j < 3; j++ {
			m[i*3+j] = 3*rand.Float32() - 1.5
		}
	}
	vv := m.mul(v)
	mm := m.inv()
	actual := mm.mul(vv)
	expected := v
	if !actual.close(expected, 3) {
		t.Errorf(`want %v got %v`, expected, actual)
	}
}
