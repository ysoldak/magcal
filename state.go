package magcal

import "fmt"

var DefaultStateData = []float32{
	0, 0, 0, // offset
	1, 0, 0, // covariance matrix
	0, 1, 0,
	0, 0, 1,
}

type State struct {
	data []float32 // data, 4x3 matrix
	off  vector    // view, offsets, first row
	cov  matrix    // view, covariance, 3x3 matrix starting from 2nd row
}

func NewState(data []float32) State {
	return State{
		data: data,
		off:  data[:3],
		cov:  data[3:12],
	}
}

func DefaultState() State {
	defaultStateCopy := make([]float32, 12)
	copy(defaultStateCopy, DefaultStateData)
	return NewState(defaultStateCopy)
}

func (s State) Export() []float32 {
	data := make([]float32, 12)
	copy(data, s.data)
	return data
}

func (s State) apply(v vector) (w vector) {
	return s.cov.mul(v.sub(s.off))
}

func (s State) dump() {
	println("===")
	for i := range s.data {
		fmt.Printf("%+0.2f ", s.data[i])
		if i > 0 && (i+1)%3 == 0 {
			println()
		}
	}
	println("===")
}

func (s State) diff(ss State) State {
	result := make([]float32, 12)
	for i := range s.data {
		result[i] = s.data[i] - ss.data[i]
	}
	return NewState(result)
}
