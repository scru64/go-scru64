package scru64

import "math/rand"

// An interface to customize the initial counter value for each new `timestamp`.
//
// [Generator] calls `Renew()` to obtain the initial counter value when the
// `timestamp` field has changed since the immediately preceding ID. Types
// implementing this interface may apply their respective logic to calculate the
// initial counter value.
type CounterMode interface {
	// Returns the next initial counter value of `counterSize` bits.
	//
	// `Generator` passes the `counterSize` (from 1 to 23) and other context
	// information that may be useful for counter renewal. The returned value must
	// be within the range of `counterSize`-bit unsigned integer.
	Renew(counterSize uint8, context CounterModeRenewContext) uint32
}

// Represents the context information provided by [Generator] to
// `CounterMode.Renew`.
type CounterModeRenewContext struct {
	// The `timestamp` value for the new counter.
	Timestamp uint64

	// The `nodeId` of the generator.
	NodeId uint32
}

// Creates a new instance of the default "initialize a portion counter" mode
// with the size (in bits) of overflow guard bits.
//
// With this mode, the counter is reset to a random number for each new
// `timestamp` tick, but some specified leading bits are set to zero to reserve
// space as the counter overflow guard.
//
// Note that the random number generator employed is not cryptographically
// strong. This mode does not pay for security because a small random number is
// insecure anyway.
func NewDefaultCounterMode(overflowGuardSize uint8) CounterMode {
	return &defaultCounterMode{overflowGuardSize: overflowGuardSize}
}

// The default "initialize a portion counter" strategy.
type defaultCounterMode struct {
	overflowGuardSize uint8
}

// Returns the next initial counter value of `counterSize` bits.
func (c *defaultCounterMode) Renew(
	counterSize uint8, _ CounterModeRenewContext) uint32 {
	if c.overflowGuardSize < counterSize {
		return rand.Uint32() >> (32 + c.overflowGuardSize - counterSize)
	} else {
		return 0
	}
}
