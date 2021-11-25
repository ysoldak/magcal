package magcal

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	mc := NewDefault()
	v := vector{1, 2, 3}
	x, y, z := mc.Apply(v[0], v[1], v[2])
	w := vector{x, y, z}
	if !v.close(w, 3) {
		t.Fatalf(`want %v got %v`, v, w)
	}
}

func TestCustom(t *testing.T) {
	cal := []float32{
		0.5, 0.0, 0.0, // offset
		1.0, 0.0, 0.0, // covariance matrix
		0.0, 1.2, 0.0,
		0.0, 0.0, 1.0,
	}
	mc := New(NewState(cal), DefaultConfiguration())
	v := vector{1, 2, 3}
	expect := vector{0.5, 2.4, 3}
	ax, ay, az := mc.Apply(v[0], v[1], v[2])
	actual := vector{ax, ay, az}
	if !actual.close(expect, 3) {
		t.Fatalf(`want %v got %v`, expect, actual)
	}
}

var calOptimal = []float32{

	// fake
	0.5, -0.2, 0.1, // offset
	1.0, -0.15, 0.0, // covariance matrix
	0.0, 1.2, 0.1,
	0.0, 0.19, 1.1,

	// real from board
	// -0.13, 0.24, 0.05,
	// 0.99, -0.15, -0.05,
	// 0.19, 1.04, -0.11,
	// 0.09, 0.14, 1.10,

}

// 5 of 8 quadrants are enough to have a good result (sometimes)
// does not matter which, 5 is enough for them be on different sides of (0,0)
func TestSearch(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	// activeQuadrants := []byte{0, 1, 2, 3, 4}
	activeQuadrants := []byte{0, 1, 2, 3, 4, 5, 6, 7}

	mc := NewDefault()

	// load buffer with random uncalibrated vectors
	for !mc.buf.full() {
		v := unCalibrate(randomCalibrated(), calOptimal)
		w := mc.State.apply(v) // it be same as v here actually, since no calibration happened yet
		q := v.quadrant()
		for _, a := range activeQuadrants {
			if q == a {
				mc.buf.push(v, w)
			}
		}
	}

	// search for solution
	iter := mc.search()

	fmt.Printf("%d %+0.3f %+0.3f\r\n", iter, mc.errorTotal(), mc.errorTotal()/float32(mc.Config.BufferSize))
	mc.State.dump()                             // found solution
	diff := mc.State.diff(NewState(calOptimal)) // difference from optimal
	diff.dump()

	// buffer
	// for _, v := range mc.buf.raw {
	// 	w := mc.State.apply(v)
	// 	e := mc.error(w)
	// 	fmt.Printf("%v, %+0.3f, %+0.3f\r\n", w.string(), w.len(), e)
	// }
	// println()
	// for _, w := range mc.buf.cal {
	// 	e := mc.error(w)
	// 	fmt.Printf("%v, %+0.3f, %+0.3f\r\n", w.string(), w.len(), e)
	// }

	// check
	for i, s := range diff.data {
		if (i < 3) && abs(s) > mc.Config.Target*0.1 {
			t.Fatal("Result calibration matrix offset item differs more than 10%: ", i, s, calOptimal[i])
		} else if (i == 3 || i == 7 || i == 11) && abs(s) > 0.1 { // trace shall be quite good
			t.Fatal("Result calibration matrix trace item differs more than 0.1: ", i, s, calOptimal[i])
		} else if abs(s) > 0.5 {
			t.Fatal("Result calibration matrix non-trace item differs more than 0.5: ", i, s, calOptimal[i])
		}
	}

	// check error for each vector in buffer
	for _, v := range mc.buf.raw {
		w := mc.State.apply(v)
		e := mc.error(w)
		if e > mc.Config.Tolerance*10 {
			fmt.Printf("%v, %+0.3f, %+0.3f\r\n", w.string(), w.len(), e)
			t.Fatalf(`error larger than expected: %0.4f`, e)
		}
	}

}

func BenchmarkSearch(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	for n := 0; n < b.N; n++ {
		mc := NewDefault()
		for i := 0; i < mc.Config.BufferSize*10; i++ {
			v := unCalibrate(randomCalibrated(), calOptimal)
			w := mc.State.apply(v) // at first v and w be same, and after buffer fills they diverge
			mc.buf.push(v, w)
			if mc.buf.full() {
				mc.search()
			}
		}
	}
}

func randomCalibrated() (v vector) {
	x := 3*rand.Float32() - 1.5
	y := 3*rand.Float32() - 1.5
	z := 3*rand.Float32() - 1.5
	random := vector{x, y, z}
	v = random.norm()
	return v
}

func unCalibrate(v vector, cal []float32) (w vector) {
	m := matrix(make([]float32, 9))
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i*3+j] = cal[(i+1)*3+j]
		}
	}
	mi := m.inv()
	w = mi.mul(v)
	for i := range w {
		w[i] += cal[i]
	}
	return w
}
