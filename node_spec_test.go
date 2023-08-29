package scru64

import (
	"database/sql"
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Initializes with node ID and size pair and node spec string.
func TestConstructor(t *testing.T) {
	for _, e := range exampleNodeSpecs {
		var nodePrev Id
		var err error
		nodePrev, err = FromUint(e.nodePrev)
		assert(t, err == nil)

		var withNodePrev, withNodeId, parsed NodeSpec
		withNodePrev, err = NewNodeSpecWithNodePrev(nodePrev, e.nodeIdSize)
		assert(t, err == nil)
		assert(t, withNodePrev.NodeId() == e.nodeId)
		assert(t, withNodePrev.NodeIdSize() == e.nodeIdSize)
		if p := withNodePrev.NodePrev(); p != 0 {
			assert(t, p == nodePrev)
		}
		assert(t, withNodePrev.nodePrev == nodePrev)
		assert(t, withNodePrev.String() == e.canonical)

		withNodeId, err = NewNodeSpecWithNodeId(e.nodeId, e.nodeIdSize)
		assert(t, err == nil)
		assert(t, withNodeId.NodeId() == e.nodeId)
		assert(t, withNodeId.NodeIdSize() == e.nodeIdSize)
		assert(t, withNodeId.NodePrev() == 0)
		if strings.HasSuffix(e.specType, "_node_id") {
			assert(t, withNodeId.nodePrev == nodePrev)
			assert(t, withNodeId.String() == e.canonical)
		}

		parsed, err = ParseNodeSpec(e.nodeSpec)
		assert(t, err == nil)
		assert(t, parsed.NodeId() == e.nodeId)
		assert(t, parsed.NodeIdSize() == e.nodeIdSize)
		if p := parsed.NodePrev(); p != 0 {
			assert(t, p == nodePrev)
		}
		assert(t, parsed.nodePrev == nodePrev)
		assert(t, parsed.String() == e.canonical)
	}
}

// Fails to initialize with invalid node spec string.
func TestConstructorError(t *testing.T) {
	var cases = []string{
		"",
		"42",
		"/8",
		"42/",
		" 42/8",
		"42/8 ",
		" 42/8 ",
		"42 / 8",
		"+42/8",
		"42/+8",
		"-42/8",
		"42/-8",
		"ab/8",
		"1/2/3",
		"0/0",
		"0/24",
		"8/1",
		"1024/8",
		"0000000000001/8",
		"1/0016",
		"42/800",
	}

	for _, e := range cases {
		var x NodeSpec
		var err error
		x, err = ParseNodeSpec(e)
		assert(t, x == NodeSpec{} && err != nil)
	}
}

// Supports serialization and deserialization.
func TestNodeSpecSerDe(t *testing.T) {
	var x, y, z NodeSpec
	var buf []byte
	var err error
	for _, e := range exampleNodeSpecs {
		x, _ = NewNodeSpecWithNodePrev(Id(e.nodePrev), e.nodeIdSize)

		buf, err = json.Marshal(x)
		assert(t, string(buf) == `"`+e.canonical+`"` && err == nil)

		err = json.Unmarshal([]byte(`"`+e.canonical+`"`), &y)
		assert(t, x == y && err == nil)

		err = z.Scan(e.nodeSpec)
		assert(t, x == z && err == nil)
	}
}

// Ensures compliance with interfaces.
func TestNodeSpecInterfaces(t *testing.T) {
	var x NodeSpec
	var _ fmt.Stringer = x
	var _ encoding.TextUnmarshaler = &x
	var _ encoding.TextMarshaler = x
	var _ sql.Scanner = &x
}
