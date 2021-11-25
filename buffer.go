package magcal

type buffer struct {
	size      int      // maximum number of vectors in buffer
	raw       []vector // raw vectors
	cal       []vector // calibrated vectors, optimise gradient calculations
	quadrants [8][]int // quadrants of raw vectors, oldest vector in most populated quadrant is replaced with new vector (quadrants may not match)
	qmax      byte     // index of quadrant with most vectors
}

func (b *buffer) full() bool {
	return len(b.raw) == b.size
}

func (b *buffer) push(v vector, w vector) {
	if !b.full() {
		b.append(v, w)
	} else {
		b.replace(v, w)
	}
}

func (b *buffer) append(v vector, w vector) {
	vq := v.quadrant()

	b.raw = append(b.raw, v)
	b.cal = append(b.cal, w)
	b.quadrants[vq] = append(b.quadrants[vq], len(b.raw)-1) // add to quadrant

	if len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
		b.qmax = vq
	}
}

func (b *buffer) replace(v vector, w vector) {
	vq := v.quadrant()

	i := b.quadrants[b.qmax][0]                   // get oldest index in quadrant
	b.quadrants[b.qmax] = b.quadrants[b.qmax][1:] // remove it from quadrant

	b.raw[i] = v
	b.cal[i] = w
	b.quadrants[vq] = append(b.quadrants[vq], i) // add to quadrant

	if len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
		b.qmax = vq
	}
}
