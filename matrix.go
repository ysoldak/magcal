package magcal

type matrix []float32 // we expect size 9 = 3x3

func (m matrix) mul(v vector) (w vector) {
	w = vector{0, 0, 0}
	for i := range v {
		w[i] = m[i*3]*v[0] + m[i*3+1]*v[1] + m[i*3+2]*v[2]
	}
	return w
}
