package magcal

type buffer struct {
	size      int
	data      []vector
	quadrants [8]([]*vector)
	qmax      byte
}

func (b *buffer) full() bool {
	return len(b.data) == b.size
}

func (b *buffer) push(v vector) {
	vq := v.quadrant()
	if !b.full() {
		b.data = append(b.data, v)
		b.quadrants[vq] = append(b.quadrants[vq], &v)
		if len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
			b.qmax = vq
		}
		return
	}
	rv := b.quadrants[b.qmax][0]
	b.quadrants[b.qmax] = b.quadrants[b.qmax][1:] // remove from quadrant
	for ci, cv := range b.data {                  // replace in data
		if &cv == rv {
			b.data[ci] = v
			break
		}
	}
	b.quadrants[vq] = append(b.quadrants[vq], &v) // add to new quadrant

	// adjust qmax if needed
	if vq != b.qmax && len(b.quadrants[vq]) > len(b.quadrants[b.qmax]) {
		b.qmax = vq
	}

}
