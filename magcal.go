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
	State State

	// internal
	config  Configuration
	active  bool
	working bool
	buf     buffer
	target2 float32 // target squared, for speed
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
		config:  config,
		active:  false,
		working: false,
		buf:     buffer{size: config.BufferSize},
		target2: config.Target * config.Target,
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
	if mc.error(w) < mc.config.Tolerance { // small error
		return
	}
	mc.buf.push(v, w)
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
				step := s * mc.config.Step
				if i < 3 { // offset must be adjusted with larger steps
					step *= mc.config.Target
				}
				for {
					if !mc.isGoodChange(i, step) {
						break
					}
					mc.adjustBuffer(i, step)
					newErr := mc.errorTotal()
					if newErr >= curErr {
						mc.adjustBuffer(i, -step) // revert
						break
					}
					// if s > 0 {
					// 	print("+")
					// } else {
					// 	print("-")
					// }
					mc.State.data[i] += step
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
	return abs(v.lenSq() - mc.target2)
	// return abs(v.len() - mc.Config.Target)
}

func (mc *MagCal) errorTotal() float32 {
	sum := float32(0)
	for _, w := range mc.buf.cal {
		sum += mc.error(w)
		if mc.config.Throttle > 0 {
			time.Sleep(mc.config.Throttle / time.Duration(mc.buf.size))
		}
	}
	return sum
}

// isGoodChange verifies is it OK to adjust
// "i"th calibration parameter by "step"
func (mc *MagCal) isGoodChange(i int, step float32) bool {
	if i < 3 { // offset is always fine
		return true
	}
	val := abs(mc.State.data[i] + step)
	if i == 3 || i == 7 || i == 11 { // covariance matrix
		return 0.5 < val && val < 2
	}
	x := abs(mc.State.data[3])
	y := abs(mc.State.data[7])
	z := abs(mc.State.data[11])
	return val*5 < x && val*5 < y && val*5 < z
}

// adjustBuffer recalculates calibrated vectors in buffer
// given "idx"th calibration parameter changed by "step"
// calVec[i] = cov[i][0]*(raw[i]-off[i]) + ... + cov[i][2]*(raw[i]-off[i])
func (mc *MagCal) adjustBuffer(idx int, step float32) {
	for n, w := range mc.buf.cal {
		if idx < 3 {
			// offset value changed
			// each coord shall be adjusted by step multiplied by respective cov value
			for i := 0; i < 3; i++ {
				w[i] -= mc.State.cov[i*3+idx] * step
			}
		} else {
			// cov value changed
			// enough to adjust one coord
			i := idx/3 - 1 // row
			j := idx % 3   // column
			w[i] += step * (mc.buf.raw[n][j] - mc.State.off[j])
		}
	}
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

func abs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}
