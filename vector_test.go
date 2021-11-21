package magcal

import (
	"fmt"
	"math"
)

// Methods here are only for tests

func (a vector) len() float32 {
	return float32(math.Sqrt(float64(a.lenSq())))
}

func (a vector) norm() (w vector) {
	len := a.len()
	return vector{a[0] / len, a[1] / len, a[2] / len}
}

func (v vector) string() string {
	return fmt.Sprintf("[%+0.3f,%+0.3f,%+0.3f]", v[0], v[1], v[2])
}

func (a vector) close(b vector, precision int) bool {
	mul := float32(10 * precision)
	for i := range a {
		if int(a[i]*mul) != int(b[i]*mul) {
			return false
		}
	}
	return true
}
