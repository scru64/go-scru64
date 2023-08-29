package scru64

import (
	"testing"
	"time"
)

func assertConsecutive(t *testing.T, first Id, second Id) {
	assert(t, first < second)
	if first.Timestamp() == second.Timestamp() {
		assert(t, first.NodeCtr()+1 == second.NodeCtr())
	} else {
		assert(t, first.Timestamp()+1 == second.Timestamp())
	}
}

// Normally generates monotonic IDs or resets state upon significant rollback.
func TestGenerateOrReset(t *testing.T) {
	const nLoops = 64
	const allowance uint64 = 10_000

	for _, e := range exampleNodeSpecs {
		counterSize := 24 - e.nodeIdSize
		nodeSpec, _ := NewNodeSpecWithNodeId(e.nodeId, e.nodeIdSize)
		g := NewGenerator(nodeSpec)

		// happy path
		var ts uint64 = 1_577_836_800_000 // 2020-01-01
		var prev, curr Id
		prev = g.GenerateOrResetCore(ts, allowance)
		for i := 0; i < nLoops; i++ {
			ts += 16
			curr = g.GenerateOrResetCore(ts, allowance)
			assertConsecutive(t, prev, curr)
			assert(t, (curr.Timestamp()-(ts>>8)) < (allowance>>8))
			assert(t, (curr.NodeCtr()>>counterSize) == e.nodeId)

			prev = curr
		}

		// keep monotonic order under mildly decreasing timestamps
		ts += allowance * 16
		prev = g.GenerateOrResetCore(ts, allowance)
		for i := 0; i < nLoops; i++ {
			ts -= 16
			curr = g.GenerateOrResetCore(ts, allowance)
			assertConsecutive(t, prev, curr)
			assert(t, (curr.Timestamp()-(ts>>8)) < (allowance>>8))
			assert(t, (curr.NodeCtr()>>counterSize) == e.nodeId)

			prev = curr
		}

		// reset state with significantly decreasing timestamps
		ts += allowance * 16
		prev = g.GenerateOrResetCore(ts, allowance)
		for i := 0; i < nLoops; i++ {
			ts -= allowance + 0x100
			curr = g.GenerateOrResetCore(ts, allowance)
			assert(t, prev > curr)
			assert(t, (curr.Timestamp()-(ts>>8)) < (allowance>>8))
			assert(t, (curr.NodeCtr()>>counterSize) == e.nodeId)

			prev = curr
		}
	}
}

// Normally generates monotonic IDs or aborts upon significant rollback.
func TestGenerateOrAbort(t *testing.T) {
	const nLoops = 64
	const allowance uint64 = 10_000

	for _, e := range exampleNodeSpecs {
		counterSize := 24 - e.nodeIdSize
		nodeSpec, _ := NewNodeSpecWithNodeId(e.nodeId, e.nodeIdSize)
		g := NewGenerator(nodeSpec)

		// happy path
		var ts uint64 = 1_577_836_800_000 // 2020-01-01
		var prev, curr Id
		var err error
		prev, err = g.GenerateOrAbortCore(ts, allowance)
		assert(t, err == nil)
		for i := 0; i < nLoops; i++ {
			ts += 16
			curr, err = g.GenerateOrAbortCore(ts, allowance)
			assert(t, err == nil)
			assertConsecutive(t, prev, curr)
			assert(t, (curr.Timestamp()-(ts>>8)) < (allowance>>8))
			assert(t, (curr.NodeCtr()>>counterSize) == e.nodeId)

			prev = curr
		}

		// keep monotonic order under mildly decreasing timestamps
		ts += allowance * 16
		prev, err = g.GenerateOrAbortCore(ts, allowance)
		assert(t, err == nil)
		for i := 0; i < nLoops; i++ {
			ts -= 16
			curr, err = g.GenerateOrAbortCore(ts, allowance)
			assert(t, err == nil)
			assertConsecutive(t, prev, curr)
			assert(t, (curr.Timestamp()-(ts>>8)) < (allowance>>8))
			assert(t, (curr.NodeCtr()>>counterSize) == e.nodeId)

			prev = curr
		}

		// abort with significantly decreasing timestamps
		ts += allowance * 16
		_, err = g.GenerateOrAbortCore(ts, allowance)
		assert(t, err == nil)
		ts -= allowance + 0x100
		for i := 0; i < nLoops; i++ {
			ts -= 16
			_, err = g.GenerateOrAbortCore(ts, allowance)
			assert(t, err == ErrClockRollback)
		}
	}
}

// Embeds up-to-date timestamp.
func TestClockIntegration(t *testing.T) {
	for _, e := range exampleNodeSpecs {
		nodeSpec, _ := NewNodeSpecWithNodeId(e.nodeId, e.nodeIdSize)
		g := NewGenerator(nodeSpec)

		var tsNow uint64
		var x Id
		var err error

		tsNow = uint64(time.Now().UnixMilli() >> 8)
		x, err = g.Generate()
		assert(t, err == nil)
		assert(t, x.Timestamp()-tsNow <= 1)

		tsNow = uint64(time.Now().UnixMilli() >> 8)
		x = g.GenerateOrReset()
		assert(t, x.Timestamp()-tsNow <= 1)

		tsNow = uint64(time.Now().UnixMilli() >> 8)
		x = g.GenerateOrSleep()
		assert(t, x.Timestamp()-tsNow <= 1)
	}
}
