package magcal

type vector []float32 // we expect size 3

func (a vector) lenSq() float32 {
	result := float32(0)
	for _, d := range a {
		result += d * d
	}
	return result
}

func (a vector) sub(b vector) (c vector) {
	return vector{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
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
