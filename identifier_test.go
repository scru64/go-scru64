package scru64

import (
	"database/sql"
	"encoding"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

// Supports equality comparison.
func TestEq(t *testing.T) {
	var prev, curr, twin Id
	prev, _ = FromUint(exampleIds[len(exampleIds)-1].num)
	for _, e := range exampleIds {
		curr, _ = FromUint(e.num)
		twin, _ = FromUint(e.num)

		assert(t, curr == twin)
		assert(t, twin == curr)
		assert(t, curr.Num() == twin.Num())
		assert(t, curr.String() == twin.String())
		assert(t, curr.Timestamp() == twin.Timestamp())
		assert(t, curr.NodeCtr() == twin.NodeCtr())

		assert(t, curr != prev)
		assert(t, prev != curr)
		assert(t, curr.Num() != prev.Num())
		assert(t, curr.String() != prev.String())
		assert(t, (curr.Timestamp() != prev.Timestamp()) ||
			(curr.NodeCtr() != prev.NodeCtr()))

		prev = curr
	}
}

// Supports ordering comparison.
func TestOrd(t *testing.T) {
	cases := append([]exampleId{}, exampleIds...)
	sort.Slice(cases, func(i, j int) bool { return cases[i].num < cases[j].num })

	var prev, curr Id
	prev, _ = FromUint(cases[0].num)
	for _, e := range cases[1:] {
		curr, _ = FromUint(e.num)

		assert(t, prev < curr)
		assert(t, prev <= curr)

		assert(t, curr > prev)
		assert(t, curr >= prev)

		assert(t, prev.Num() < curr.Num())
		assert(t, prev.String() < curr.String())

		prev = curr
	}
}

// Converts to various types.
func TestConvertTo(t *testing.T) {
	var x Id
	var buf []byte
	var err error

	for _, e := range exampleIds {
		x, _ = FromUint(e.num)

		assert(t, x.Num() == e.num)
		assert(t, uint64(x) == e.num)
		assert(t, int64(x) == int64(e.num))
		assert(t, x.String() == e.text)
		assert(t, fmt.Sprint(x) == e.text)
		buf, err = x.MarshalText()
		assert(t, string(buf) == e.text && err == nil)
		assert(t, x.Timestamp() == e.timestamp)
		assert(t, x.NodeCtr() == e.nodeCtr)
	}
}

// Converts from various types.
func TestConvertFrom(t *testing.T) {
	var x, y Id
	var err error

	for _, e := range exampleIds {
		x, _ = FromUint(e.num)

		y, err = Parse(e.text)
		assert(t, x == y && err == nil)
		y, err = Parse(strings.ToUpper(e.text))
		assert(t, x == y && err == nil)
		y = 0
		err = y.UnmarshalText([]byte(e.text))
		assert(t, x == y && err == nil)
		y = 0
		err = y.UnmarshalText([]byte(strings.ToUpper(e.text)))
		assert(t, x == y && err == nil)
		y = 0
		err = y.Scan(int64(e.num))
		assert(t, x == y && err == nil)
		y = 0
		err = y.Scan(e.text)
		assert(t, x == y && err == nil)
		y = 0
		err = y.Scan(strings.ToUpper(e.text))
		assert(t, x == y && err == nil)

		y, err = FromParts(e.timestamp, e.nodeCtr)
		assert(t, x == y && err == nil)
	}
}

// Rejects integer out of valid range.
func TestFromIntError(t *testing.T) {
	var x Id
	var err error

	x, err = FromUint(4738381338321616896)
	assert(t, x == 0 && err != nil)
	x, err = FromUint(0xffff_ffff_ffff_ffff)
	assert(t, x == 0 && err != nil)
}

// Fails to parse invalid textual representations.
func TestParseError(t *testing.T) {
	cases := []string{
		"",
		" 0u3wrp5g81jx",
		"0u3wrp5g81jy ",
		" 0u3wrp5g81jz ",
		"+0u3wrp5g81k0",
		"-0u3wrp5g81k1",
		"+u3wrp5q7ta5",
		"-u3wrp5q7ta6",
		"0u3w_p5q7ta7",
		"0u3wrp5-7ta8",
		"0u3wrp5q7t 9",
	}

	var x Id
	var err error
	for _, e := range cases {
		x, err = Parse(e)
		assert(t, x == 0 && err != nil)
	}
}

// Rejects `MAX + 1` even if passed as pair of fields.
func TestFromPartsError(t *testing.T) {
	var max uint64 = 4738381338321616895
	var x Id
	var err error
	x, err = FromParts(max>>24, uint32(max&0xff_ffff)+1)
	assert(t, x == 0 && err != nil)
}

// Supports serialization and deserialization.
func TestSerDe(t *testing.T) {
	var x, y Id
	var buf []byte
	var err error
	for _, e := range exampleIds {
		x, _ = FromUint(e.num)

		buf, err = json.Marshal(x)
		assert(t, string(buf) == `"`+e.text+`"` && err == nil)

		err = json.Unmarshal([]byte(`"`+e.text+`"`), &y)
		assert(t, x == y && err == nil)
	}
}

// Ensures compliance with interfaces.
func TestInterfaces(t *testing.T) {
	var x Id
	var _ fmt.Stringer = x
	var _ encoding.TextUnmarshaler = &x
	var _ encoding.TextMarshaler = x
	var _ sql.Scanner = &x
}
