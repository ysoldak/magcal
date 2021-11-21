package magcal

import (
	"fmt"
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

// Helper methods

func (m matrix) string() string {
	result := ""
	for i := 0; i < 3; i++ {
		result += fmt.Sprintf("[%0.3f,%0.3f,%0.3f]\r\n", m[i*3], m[i*3+1], m[i*3+2])
	}
	return result
}

// Calculate inverse matrix
// https://www.wikihow.com/Find-the-Inverse-of-a-3x3-Matrix

func (m matrix) det() float32 {
	d1 := m[0] * m.detSub(0, 0)
	d2 := m[1] * m.detSub(0, 1)
	d3 := m[2] * m.detSub(0, 2)
	return d1 - d2 + d3
}

func (m matrix) detSub(x, y int) float32 {
	i1, i2, j1, j2 := 0, 0, 0, 0
	if x == 0 {
		i1 = 1
		i2 = 2
	}
	if x == 1 {
		i1 = 0
		i2 = 2
	}
	if x == 2 {
		i1 = 0
		i2 = 1
	}
	if y == 0 {
		j1 = 1
		j2 = 2
	}
	if y == 1 {
		j1 = 0
		j2 = 2
	}
	if y == 2 {
		j1 = 0
		j2 = 1
	}
	return m[i1*3+j1]*m[i2*3+j2] - m[i1*3+j2]*m[i2*3+j1]
}

func (m matrix) trans() (w matrix) {
	w = make([]float32, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w[j*3+i] = m[i*3+j]
		}
	}
	return w
}

func (m matrix) inv() (w matrix) {
	d := m.det()
	t := m.trans()
	w = make([]float32, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w[i*3+j] = t.detSub(i, j) / d
			if (i == 1 || j == 1) && (i*j != 1) {
				w[i*3+j] *= -1
			}
		}
	}
	return w
}

func (m matrix) close(w matrix, precision int) bool {
	mul := float32(10 * precision)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if int(m[i*3+j]*mul) != int(w[i*3+j]*mul) {
				return false
			}
		}
	}
	return true
}
