package magcal

import (
	"time"
)

const (
	DefaultTarget     = 1.0
	DefautlTolerance  = 0.01
	DefaultStep       = 0.01
	DefaultBufferSize = 128
)

type MagCal struct {
	State  State
	Config Configuration

	// internal
	active  bool
	working bool
	buf     buffer
}

type Configuration struct {
	Target     float32       // vector lengths after calibration must be as close to this as possible
	Step       float32       // step size on covariance matrix trace, offset step is target*step
	Tolerance  float32       // abs(len^2-target)
	BufferSize int           // number of vectors in cache to run gradient descent on
	Throttle   time.Duration // increase to let other goroutines to work more by throttling search goroutine
}

func DefaultConfiguration() Configuration {
	return Configuration{
		Target:     DefaultTarget,
		Step:       DefaultStep,
		Tolerance:  DefautlTolerance,
		BufferSize: DefaultBufferSize,
		Throttle:   0,
	}
}

func NewDefault() *MagCal {
	return New(DefaultState(), DefaultConfiguration())
}

func New(state State, config Configuration) *MagCal {
	return &MagCal{
		State:   state,
		Config:  config,
		active:  false,
		working: false,
		buf:     buffer{size: config.BufferSize},
	}
}

func (mc *MagCal) Start() {
	mc.active = true
}

func (mc *MagCal) Stop() {
	mc.active = false
}

func (mc *MagCal) Apply(x, y, z float32) (xx, yy, zz float32) {
	v := vector{x, y, z}
	w := mc.State.apply(v) // calibrated
	xx, yy, zz = w[0], w[1], w[2]
	if !mc.active {
		return
	}
	if mc.working {
		return
	}
	if mc.error(w) < mc.Config.Tolerance { // small error
		return
	}
	mc.buf.push(v)
	if !mc.buf.full() {
		return
	}
	go mc.search()
	return
}

func (mc *MagCal) search() int {
	mc.working = true

	// offests, cov trace, rest
	indices := []int{0, 1, 2, 3, 7, 11, 4, 5, 6, 8, 9, 10}

	curErr := mc.errorTotal()
	iter := 0
	improved := true
	// println(curErr)
	// println()
	for improved {
		iter++
		improved = false
		// print("#")
		for _, i := range indices {
			// print("*")
			for _, s := range [2]float32{-1, 1} {
				step := s * mc.Config.Step
				for {
					mc.State.data[i] += step
					if !mc.isGoodChange(i) {
						mc.State.data[i] -= step // revert
						break
					}
					newErr := mc.errorTotal()
					if newErr >= curErr {
						mc.State.data[i] -= step // revert
						break
					}
					// if s > 0 {
					// 	print("+")
					// } else {
					// 	print("-")
					// }
					curErr = newErr
					improved = true
				}
			}
		}
		// println()
		// mc.State.dump()
	}
	mc.working = false
	return iter
}

// --- utils ---

func (mc *MagCal) error(v vector) float32 {
	return abs(v.lenSq() - mc.Config.Target)
	// return abs(v.len() - DefaultTarget)
}

func (mc *MagCal) errorTotal() float32 {
	sum := float32(0)
	for _, v := range mc.buf.data {
		w := mc.State.apply(v)
		sum += mc.error(w)
		if mc.Config.Throttle > 0 {
			time.Sleep(mc.Config.Throttle / time.Duration(mc.buf.size))
		}
	}
	return sum
}

// trace values (of cov matrix):
// none of them can be more than twice bigger than other trace values
//
// non-trace values:
// shall be at least 5 times smaller than any value in trace
// func (mc *MagCal) isGoodCal() bool {
// 	x := abs(mc.state[3])
// 	y := abs(mc.state[7])
// 	z := abs(mc.state[11])
// 	if x < 0.5 || y < 0.5 || z < 0.5 || x > 2 || y > 2 || z > 2 {
// 		return false
// 	}
// 	if x > y*2 || x > z*2 || y > x*2 || y > z*2 || z > x*2 || z > y*2 {
// 		return false
// 	}
// 	for i := 4; i < 11; i++ {
// 		if i == 7 { // ignore y
// 			continue
// 		}
// 		v := abs(mc.state[i]) * 5
// 		if v > x || v > y || v > z {
// 			return false
// 		}
// 	}
// 	return true
// }

func (mc *MagCal) isGoodChange(i int) bool {
	if i < 3 {
		return true
	}
	val := abs(mc.State.data[i])
	if i == 3 || i == 7 || i == 11 {
		return 0.5 < val && val < 2
	}
	x := abs(mc.State.data[3])
	y := abs(mc.State.data[7])
	z := abs(mc.State.data[11])
	return val*5 < x && val*5 < y && val*5 < z
}

func abs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}
