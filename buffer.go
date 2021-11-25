package magcal

type buffer struct {
	size      int            // maximum number of vectors in buffer
	raw       []vector       // raw vectors
	cal       []vector       // calibrated vectors, optimise gradient calculations
	quadrants [8]([]*vector) // quadrants of raw vectors, used to determine which existing vector shall go on push
	qmax      byte           // index of quadrant with most vectors
}

func (b *buffer) full() bool {
	return len(b.raw) == b.size
}

func (b *buffer) push(v vector, w vector) {
	vq := v.quadrant()
	if !b.full() {
		b.raw = append(b.raw, v)
		b.cal = append(b.cal, w)
		b.quadrants[vq] = append(b.quadrants[vq], &v)
		if len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
			b.qmax = vq
		}
		return
	}
	rv := b.quadrants[b.qmax][0]
	b.quadrants[b.qmax] = b.quadrants[b.qmax][1:] // remove from quadrant
	for ci, cv := range b.raw {                   // replace in raw and cal
		if &cv == rv {
			b.raw[ci] = v
			b.cal[ci] = w
			break
		}
	}
	b.quadrants[vq] = append(b.quadrants[vq], &v) // add to new quadrant

	// adjust qmax if needed
	if vq != b.qmax && len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
		b.qmax = vq
	}

}
