package scru64

import "fmt"

// The maximum valid value of the `timestamp` field.
const maxTimestamp uint64 = uint64(MaxId) >> nodeCtrSize

// The maximum valid value of the combined `nodeCtr` field.
const maxNodeCtr uint32 = (1 << nodeCtrSize) - 1

// Digit characters used in the Base36 notation.
var digits = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

// An O(1) map from ASCII code points to Base36 digit values.
var decodeMap = [256]byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x01, 0x02, 0x03,
	0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
	0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
	0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// Represents a SCRU64 ID.
type Id uint64

// The minimum valid SCRU64 ID value (i.e., `000000000000`).
const MinId Id = 0

// The maximum valid SCRU64 ID value (i.e., `zzzzzzzzzzzz`).
const MaxId Id = 4738381338321616895

// Ensures that the receiver is a valid SCRU64 ID value, or panics if not.
func (n Id) verify() {
	if n > MaxId {
		panic("method call on invalid receiver")
	}
}

// Creates a value from a 64-bit unsigned integer.
//
// This function returns a non-nil error if the argument is larger than
// `36^12 - 1`.
func FromUint(value uint64) (Id, error) {
	if value > uint64(MaxId) {
		return Id(0), newRangeError(fmt.Errorf("`%T` out of range: %[1]v", value))
	}
	return Id(value), nil
}

// Returns the integer representation.
func (n Id) Uint() uint64 {
	n.verify()
	return uint64(n)
}

// Creates a value from the `timestamp` and the combined `nodeCtr` field value.
//
// This function returns a non-nil error if any argument is larger than their
// respective maximum value (`36^12 / 2^24 - 1` and `2^24 - 1`, respectively).
func FromParts(timestamp uint64, nodeCtr uint32) (Id, error) {
	if timestamp > maxTimestamp {
		return Id(0), fmt.Errorf(
			"scru64.Id: could not create SCRU64 ID from parts: `timestamp` out of range")
	} else if nodeCtr > maxNodeCtr {
		return Id(0), fmt.Errorf(
			"scru64.Id: could not create SCRU64 ID from parts: `nodeCtr` out of range")
	}
	// no further check is necessary because `MAX_SCRU64_INT` happens to equal
	// `MAX_TIMESTAMP << 24 | MAX_NODE_CTR`
	return Id(timestamp<<nodeCtrSize) | Id(nodeCtr), nil
}

// A convenient panicking version of [FromParts] for internal uses.
func mustFromParts(timestamp uint64, nodeCtr uint32) Id {
	n, err := FromParts(timestamp, nodeCtr)
	if err != nil {
		panic(err)
	}
	return n
}

// Returns the `timestamp` field value.
func (n Id) Timestamp() uint64 {
	n.verify()
	return uint64(n) >> nodeCtrSize
}

// Returns the `nodeId` and `counter` field values combined as a single integer.
func (n Id) NodeCtr() uint32 {
	n.verify()
	return uint32(n) & maxNodeCtr
}

// Creates a value from a 12-digit string representation.
//
// This function returns a non-nil error if the argument is not a valid string
// representation.
func Parse(value string) (Id, error) {
	var n Id
	return n, n.UnmarshalText([]byte(value))
}

// Returns the 12-digit canonical string representation.
func (n Id) String() string {
	v, _ := n.MarshalText()
	return string(v)
}

// See encoding.TextUnmarshaler
func (n *Id) UnmarshalText(text []byte) error {
	if n == nil {
		return fmt.Errorf("scru64.Id: method call on nil receiver")
	}
	if len(text) != 12 {
		return newParseError(fmt.Errorf(
			"invalid length: %d bytes (expected 12)", len(text)))
	}

	var v Id = 0
	for i, e := range text {
		if decodeMap[e] < 36 {
			v = v*36 + Id(decodeMap[e])
		} else if e < 0x80 {
			return newParseError(fmt.Errorf("invalid digit %q at %d", e, i))
		} else {
			return newParseError(fmt.Errorf("found non-ASCII digit at %d", i))
		}
	}
	*n = v
	return nil
}

// See encoding.TextMarshaler
func (n Id) MarshalText() (text []byte, err error) {
	n.verify()
	text = make([]byte, 12)
	for i := len(text) - 1; i >= 0; i-- {
		text[i] = digits[n%36]
		n /= 36
	}
	return
}

// See database/sql.Scanner
func (n *Id) Scan(src any) error {
	if n == nil {
		return fmt.Errorf("scru64.Id: method call on nil receiver")
	}
	switch src := src.(type) {
	case int64:
		if src < 0 || src > int64(MaxId) {
			return newRangeError(fmt.Errorf("`%T` out of range: %[1]v", src))
		}
		*n = Id(src)
		return nil
	case string:
		return n.UnmarshalText([]byte(src))
	default:
		return fmt.Errorf("scru64.Id: Scan: unsupported type conversion")
	}
}

// Wraps a raw range error to construct a unified error message.
func newRangeError(err error) error {
	return fmt.Errorf(
		"scru64.Id: could not convert integer to SCRU64 ID: %w", err)
}

// Wraps a raw parsing error to construct a unified error message.
func newParseError(err error) error {
	return fmt.Errorf("scru64.Id: could not parse string as SCRU64 ID: %w", err)
}
