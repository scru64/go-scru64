package scru64

import (
	"math"
	"testing"
)

// `DefaultCounterMode` returns random numbers, setting the leading guard bits
// to zero.
//
// This case includes statistical tests for the random number generator and thus
// may fail at a certain low probability.
func TestDefaultCounterMode(t *testing.T) {
	const N = 4096

	// set margin based on binom dist 99.999999% confidence interval
	var margin float64 = 5.730729 * math.Sqrt(0.5*0.5/N)

	context := CounterModeRenewContext{Timestamp: 0x0123_4567_89ab, NodeId: 0}
	for counterSize := uint8(1); counterSize < nodeCtrSize; counterSize++ {
		for overflowGuardSize := uint8(0); overflowGuardSize < nodeCtrSize; overflowGuardSize++ {
			// count number of set bits by bit position (from LSB to MSB)
			var countsByPos [nodeCtrSize]uint32

			var c CounterMode = NewDefaultCounterMode(overflowGuardSize)
			for i := 0; i < N; i++ {
				var n uint32 = c.Renew(counterSize, context)
				for j := range countsByPos {
					countsByPos[j] += n & 1
					n >>= 1
				}
				assert(t, n == 0)
			}

			var filled uint8 = 0
			if overflowGuardSize < counterSize {
				filled = counterSize - overflowGuardSize
			}
			for _, e := range countsByPos[:filled] {
				assert(t, math.Abs(float64(e)/N-0.5) < margin)
			}
			for _, e := range countsByPos[filled:] {
				assert(t, e == 0)
			}
		}
	}
}
