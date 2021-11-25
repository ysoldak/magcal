package magcal

import "testing"

func TestFull(t *testing.T) {
	b := buffer{
		size: 10,
	}
	for i := 0; i < 9; i++ {
		b.push(vector{1, 2, 3}, vector{1, 2, 3})
		if b.full() || len(b.raw) != i+1 {
			t.Fatal("shall not be full at", i)
		}
	}
	b.push(vector{1, 2, 3}, vector{1, 2, 3})
	if !b.full() || len(b.raw) != 10 {
		t.Fatal("shall be full at 10")
	}

	b.push(vector{1, 2, 3}, vector{1, 2, 3})
	if !b.full() || len(b.raw) != 10 {
		t.Fatal("shall be full at 11")
	}

}
