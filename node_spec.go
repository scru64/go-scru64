package scru64

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

// The unified format string for `nodeIdSize` range errors.
const fmtNodeIdSizeError = "scru64.NodeSpec: `nodeIdSize` (%v) must range from 1 to 23"

// A lazy initialization holder of the compiled regular expression representing
// the node spec syntax.
var reNodeSpecHolder struct {
	once  sync.Once
	inner *regexp.Regexp
}

// Represents a node configuration specifier used to build a [Generator].
//
// A `NodeSpec` is usually expressed as a node spec string, which starts with a
// decimal `nodeId`, a hexadecimal `nodeId` prefixed by "0x", or a 12-digit
// `nodePrev` SCRU64 ID value, followed by a slash and a decimal `nodeIdSize`
// value ranging from 1 to 23 (e.g., "42/8", "0xb00/12", "0u2r85hm2pt3/16"). The
// first and second forms create a fresh new generator with the given `nodeId`,
// while the third form constructs one that generates subsequent SCRU64 IDs to
// the `nodePrev`. See also [the usage notes] in the SCRU64 spec for tips and
// techniques to design node configurations.
//
// [the usage notes]: https://github.com/scru64/spec#informative-usage-notes
type NodeSpec struct {
	nodePrev   Id
	nodeIdSize uint8
}

// Heuristically ensures that the receiver is initialized by valid constructors,
// or panics if not.
func (n NodeSpec) verify() {
	if n.nodeIdSize == 0 {
		panic("method call on invalid receiver")
	}
}

// Returns the `nodeIdSize` value.
func (n NodeSpec) NodeIdSize() uint8 {
	n.verify()
	return n.nodeIdSize
}

// Creates an instance of [NodeSpec] with `nodePrev` and `nodeIdSize` values.
//
// This function returns a non-nil error if the `nodeIdSize` is zero or greater
// than 23.
func NewNodeSpecWithNodePrev(nodePrev Id, nodeIdSize uint8) (NodeSpec, error) {
	if 0 < nodeIdSize && nodeIdSize < nodeCtrSize {
		return NodeSpec{nodePrev: nodePrev, nodeIdSize: nodeIdSize}, nil
	} else {
		return NodeSpec{}, fmt.Errorf(fmtNodeIdSizeError, nodeIdSize)
	}
}

// Returns the `nodePrev` value if the `NodeSpec` is constructed with one or the
// zero value (`scru64.Id(0)`) otherwise.
func (n NodeSpec) NodePrev() Id {
	n.verify()
	if n.nodePrev.Timestamp() > 0 {
		return n.nodePrev
	} else {
		return Id(0)
	}
}

// Creates an instance of [NodeSpec] with `nodeId` and `nodeIdSize` values.
//
// This function returns a non-nil error if the `nodeIdSize` is zero or greater
// than 23 or if the `nodeId` does not fit in `nodeIdSize` bits.
func NewNodeSpecWithNodeId(nodeId uint32, nodeIdSize uint8) (NodeSpec, error) {
	if 0 < nodeIdSize && nodeIdSize < nodeCtrSize {
		if nodeId < (1 << nodeIdSize) {
			counterSize := nodeCtrSize - nodeIdSize
			return NodeSpec{
				nodePrev:   mustFromParts(0, nodeId<<counterSize),
				nodeIdSize: nodeIdSize,
			}, nil
		} else {
			return NodeSpec{}, fmt.Errorf(
				"scru64.NodeSpec: `nodeId` (%v) must fit in `nodeIdSize` (%v) bits",
				nodeId, nodeIdSize)
		}
	} else {
		return NodeSpec{}, fmt.Errorf(fmtNodeIdSizeError, nodeIdSize)
	}
}

// Returns the `nodeId` value given at instance creation or encoded in the
// `nodePrev` value.
func (n NodeSpec) NodeId() uint32 {
	counterSize := nodeCtrSize - n.NodeIdSize()
	return n.nodePrev.NodeCtr() >> counterSize
}

// Creates an instance of [NodeSpec] from a node spec string.
//
// This function returns a non-nil error if if an invalid node spec string is
// passed.
func ParseNodeSpec(value string) (NodeSpec, error) {
	var n NodeSpec
	return n, n.UnmarshalText([]byte(value))
}

// Returns the node spec string representation.
func (n NodeSpec) String() string {
	if n.NodePrev() != 0 {
		return fmt.Sprintf("%v/%v", n.NodePrev(), n.NodeIdSize())
	} else {
		return fmt.Sprintf("%v/%v", n.NodeId(), n.NodeIdSize())
	}
}

// See encoding.TextUnmarshaler
func (n *NodeSpec) UnmarshalText(text []byte) error {
	if n == nil {
		return fmt.Errorf("scru64.NodeSpec: method call on nil receiver")
	}

	reNodeSpecHolder.once.Do(func() {
		reNodeSpecHolder.inner = regexp.MustCompile(
			`^(?:([0-9A-Za-z]{12})|([0-9]{1,8})|0[Xx]([0-9A-Fa-f]{1,6}))/([0-9]{1,3})$`)
	})

	var m [][]byte = reNodeSpecHolder.inner.FindSubmatch(text)
	if m == nil {
		return fmt.Errorf(
			`scru64.NodeSpec: could not parse string as node spec (expected: e.g., "42/8", "0xb00/12", "0u2r85hm2pt3/16")`)
	}

	nodeIdSize, _ := strconv.ParseUint(string(m[4]), 10, 32)
	if nodeIdSize > 0xff {
		return fmt.Errorf(fmtNodeIdSizeError, nodeIdSize)
	}

	var result NodeSpec
	var err error
	if m[1] != nil {
		var nodePrev Id
		_ = nodePrev.UnmarshalText(m[1])
		result, err = NewNodeSpecWithNodePrev(nodePrev, uint8(nodeIdSize))
	} else if m[2] != nil {
		nodeId, _ := strconv.ParseUint(string(m[2]), 10, 32)
		result, err = NewNodeSpecWithNodeId(uint32(nodeId), uint8(nodeIdSize))
	} else if m[3] != nil {
		nodeId, _ := strconv.ParseUint(string(m[3]), 16, 32)
		result, err = NewNodeSpecWithNodeId(uint32(nodeId), uint8(nodeIdSize))
	} else {
		panic("unreachable")
	}

	if err == nil {
		*n = result
	}
	return err
}

// See encoding.TextMarshaler
func (n NodeSpec) MarshalText() (text []byte, err error) {
	return []byte(n.String()), nil
}

// See database/sql.Scanner
func (n *NodeSpec) Scan(src any) error {
	if n == nil {
		return fmt.Errorf("scru64.NodeSpec: method call on nil receiver")
	}
	switch src := src.(type) {
	case string:
		return n.UnmarshalText([]byte(src))
	default:
		return fmt.Errorf("scru64.NodeSpec: Scan: unsupported type conversion")
	}
}
