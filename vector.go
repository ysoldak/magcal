package magcal

import (
	"fmt"
	"math"
)

type vector []float32 // we expect size 3

func (a vector) len() float32 {
	return float32(math.Sqrt(float64(a.lenSq())))
}

func (a vector) lenSq() float32 {
	result := float32(0)
	for _, d := range a {
		result += d * d
	}
	return result
}

func (a vector) norm() (w vector) {
	len := a.len()
	return vector{a[0] / len, a[1] / len, a[2] / len}
}

func (a vector) sub(b vector) (c vector) {
	return vector{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
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

// quadrant returns number of quadrant this vector belongs to
// numbers are from 0b000 to 0b111
func (v vector) quadrant() byte {
	result := byte(0)
	for i := range v {
		if v[i] > 0 {
			result |= 1 << i
		}
	}
	return result
}

func (v vector) string() string {
	return fmt.Sprintf("[%+0.3f,%+0.3f,%+0.3f]", v[0], v[1], v[2])
}
