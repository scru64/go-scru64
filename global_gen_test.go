package scru64

import (
	"testing"
)

// Reads configuration from environment var.
func TestDefaultInitializer(t *testing.T) {
	t.Setenv("SCRU64_NODE_SPEC", "42/8")

	assert(t, GlobalGenerator.NodeId() == 42)
	assert(t, GlobalGenerator.NodeIdSize() == 8)
}

// Generates 100k monotonically increasing IDs.
func TestNewString(t *testing.T) {
	t.Setenv("SCRU64_NODE_SPEC", "42/8")

	var prev string = NewString()
	for i := 0; i < 100_000; i++ {
		var curr string = NewString()
		assert(t, prev < curr)
		prev = curr
	}
}
