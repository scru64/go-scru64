package scru64

import (
	"path/filepath"
	"runtime"
	"testing"
)

func assert(t *testing.T, condition bool) {
	if !condition {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Errorf("assertion failed at %v:%v", filepath.Base(file), line)
		} else {
			t.Errorf("assertion failed")
		}
	}
}

type exampleId struct {
	text      string
	num       uint64
	timestamp uint64
	nodeCtr   uint32
}

type exampleNodeSpec struct {
	nodeSpec   string
	canonical  string
	specType   string
	nodeId     uint32
	nodeIdSize uint8
	nodePrev   uint64
}

var exampleIds = []exampleId{
	{text: "000000000000", num: 0x0000000000000000, timestamp: 0, nodeCtr: 0},
	{text: "00000009zldr", num: 0x0000000000ffffff, timestamp: 0, nodeCtr: 16777215},
	{text: "zzzzzzzq0em8", num: 0x41c21cb8e0000000, timestamp: 282429536480, nodeCtr: 0},
	{text: "zzzzzzzzzzzz", num: 0x41c21cb8e0ffffff, timestamp: 282429536480, nodeCtr: 16777215},
	{text: "0u375nxqh5cq", num: 0x0186d52bbe2a635a, timestamp: 6557084606, nodeCtr: 2777946},
	{text: "0u375nxqh5cr", num: 0x0186d52bbe2a635b, timestamp: 6557084606, nodeCtr: 2777947},
	{text: "0u375nxqh5cs", num: 0x0186d52bbe2a635c, timestamp: 6557084606, nodeCtr: 2777948},
	{text: "0u375nxqh5ct", num: 0x0186d52bbe2a635d, timestamp: 6557084606, nodeCtr: 2777949},
	{text: "0u375ny0glr0", num: 0x0186d52bbf2a4a1c, timestamp: 6557084607, nodeCtr: 2771484},
	{text: "0u375ny0glr1", num: 0x0186d52bbf2a4a1d, timestamp: 6557084607, nodeCtr: 2771485},
	{text: "0u375ny0glr2", num: 0x0186d52bbf2a4a1e, timestamp: 6557084607, nodeCtr: 2771486},
	{text: "0u375ny0glr3", num: 0x0186d52bbf2a4a1f, timestamp: 6557084607, nodeCtr: 2771487},
	{text: "jdsf1we3ui4f", num: 0x2367c8dfb2e6d23f, timestamp: 152065073074, nodeCtr: 15127103},
	{text: "j0afcjyfyi98", num: 0x22b86eaad6b2f7ec, timestamp: 149123148502, nodeCtr: 11728876},
	{text: "ckzyfc271xsn", num: 0x16fc214296b29057, timestamp: 98719318678, nodeCtr: 11702359},
	{text: "t0vgc4c4b18n", num: 0x3504295badc14f07, timestamp: 227703085997, nodeCtr: 12668679},
	{text: "mwcrtcubk7bp", num: 0x29d3c7553e748515, timestamp: 179646715198, nodeCtr: 7636245},
	{text: "g9ye86pgplu7", num: 0x1dbb24363718aecf, timestamp: 127693764151, nodeCtr: 1617615},
	{text: "qmez19t9oeir", num: 0x30a122fef7cd6c83, timestamp: 208861855479, nodeCtr: 13462659},
	{text: "d81r595fq52m", num: 0x18278838f0660f2e, timestamp: 103742454000, nodeCtr: 6688558},
	{text: "v0rbps7ay8ks", num: 0x38a9e683bb4425ec, timestamp: 243368625083, nodeCtr: 4466156},
	{text: "z0jndjt42op2", num: 0x3ff596748ea77186, timestamp: 274703217806, nodeCtr: 10973574},
	{text: "f2bembkd4zrb", num: 0x1b844eb5d1aebb07, timestamp: 118183867857, nodeCtr: 11451143},
	{text: "mkg0fd5p76pp", num: 0x29391373ab449abd, timestamp: 177051235243, nodeCtr: 4496061},
}

var exampleNodeSpecs = []exampleNodeSpec{
	{nodeSpec: "0/1", canonical: "0/1", specType: "dec_node_id", nodeId: 0, nodeIdSize: 1, nodePrev: 0x0000000000000000},
	{nodeSpec: "1/1", canonical: "1/1", specType: "dec_node_id", nodeId: 1, nodeIdSize: 1, nodePrev: 0x0000000000800000},
	{nodeSpec: "0/8", canonical: "0/8", specType: "dec_node_id", nodeId: 0, nodeIdSize: 8, nodePrev: 0x0000000000000000},
	{nodeSpec: "42/8", canonical: "42/8", specType: "dec_node_id", nodeId: 42, nodeIdSize: 8, nodePrev: 0x00000000002a0000},
	{nodeSpec: "255/8", canonical: "255/8", specType: "dec_node_id", nodeId: 255, nodeIdSize: 8, nodePrev: 0x0000000000ff0000},
	{nodeSpec: "0/16", canonical: "0/16", specType: "dec_node_id", nodeId: 0, nodeIdSize: 16, nodePrev: 0x0000000000000000},
	{nodeSpec: "334/16", canonical: "334/16", specType: "dec_node_id", nodeId: 334, nodeIdSize: 16, nodePrev: 0x0000000000014e00},
	{nodeSpec: "65535/16", canonical: "65535/16", specType: "dec_node_id", nodeId: 65535, nodeIdSize: 16, nodePrev: 0x0000000000ffff00},
	{nodeSpec: "0/23", canonical: "0/23", specType: "dec_node_id", nodeId: 0, nodeIdSize: 23, nodePrev: 0x0000000000000000},
	{nodeSpec: "123456/23", canonical: "123456/23", specType: "dec_node_id", nodeId: 123456, nodeIdSize: 23, nodePrev: 0x000000000003c480},
	{nodeSpec: "8388607/23", canonical: "8388607/23", specType: "dec_node_id", nodeId: 8388607, nodeIdSize: 23, nodePrev: 0x0000000000fffffe},
	{nodeSpec: "0x0/1", canonical: "0/1", specType: "hex_node_id", nodeId: 0, nodeIdSize: 1, nodePrev: 0x0000000000000000},
	{nodeSpec: "0x1/1", canonical: "1/1", specType: "hex_node_id", nodeId: 1, nodeIdSize: 1, nodePrev: 0x0000000000800000},
	{nodeSpec: "0xb/8", canonical: "11/8", specType: "hex_node_id", nodeId: 11, nodeIdSize: 8, nodePrev: 0x00000000000b0000},
	{nodeSpec: "0x8f/8", canonical: "143/8", specType: "hex_node_id", nodeId: 143, nodeIdSize: 8, nodePrev: 0x00000000008f0000},
	{nodeSpec: "0xd7/8", canonical: "215/8", specType: "hex_node_id", nodeId: 215, nodeIdSize: 8, nodePrev: 0x0000000000d70000},
	{nodeSpec: "0xbaf/16", canonical: "2991/16", specType: "hex_node_id", nodeId: 2991, nodeIdSize: 16, nodePrev: 0x00000000000baf00},
	{nodeSpec: "0x10fa/16", canonical: "4346/16", specType: "hex_node_id", nodeId: 4346, nodeIdSize: 16, nodePrev: 0x000000000010fa00},
	{nodeSpec: "0xcc83/16", canonical: "52355/16", specType: "hex_node_id", nodeId: 52355, nodeIdSize: 16, nodePrev: 0x0000000000cc8300},
	{nodeSpec: "0xc8cd1/23", canonical: "822481/23", specType: "hex_node_id", nodeId: 822481, nodeIdSize: 23, nodePrev: 0x00000000001919a2},
	{nodeSpec: "0x26eff5/23", canonical: "2551797/23", specType: "hex_node_id", nodeId: 2551797, nodeIdSize: 23, nodePrev: 0x00000000004ddfea},
	{nodeSpec: "0x7c6bc4/23", canonical: "8154052/23", specType: "hex_node_id", nodeId: 8154052, nodeIdSize: 23, nodePrev: 0x0000000000f8d788},
	{nodeSpec: "v0rbps7ay8ks/1", canonical: "v0rbps7ay8ks/1", specType: "node_prev", nodeId: 0, nodeIdSize: 1, nodePrev: 0x38a9e683bb4425ec},
	{nodeSpec: "v0rbps7ay8ks/8", canonical: "v0rbps7ay8ks/8", specType: "node_prev", nodeId: 68, nodeIdSize: 8, nodePrev: 0x38a9e683bb4425ec},
	{nodeSpec: "v0rbps7ay8ks/16", canonical: "v0rbps7ay8ks/16", specType: "node_prev", nodeId: 17445, nodeIdSize: 16, nodePrev: 0x38a9e683bb4425ec},
	{nodeSpec: "v0rbps7ay8ks/23", canonical: "v0rbps7ay8ks/23", specType: "node_prev", nodeId: 2233078, nodeIdSize: 23, nodePrev: 0x38a9e683bb4425ec},
	{nodeSpec: "z0jndjt42op2/1", canonical: "z0jndjt42op2/1", specType: "node_prev", nodeId: 1, nodeIdSize: 1, nodePrev: 0x3ff596748ea77186},
	{nodeSpec: "z0jndjt42op2/8", canonical: "z0jndjt42op2/8", specType: "node_prev", nodeId: 167, nodeIdSize: 8, nodePrev: 0x3ff596748ea77186},
	{nodeSpec: "z0jndjt42op2/16", canonical: "z0jndjt42op2/16", specType: "node_prev", nodeId: 42865, nodeIdSize: 16, nodePrev: 0x3ff596748ea77186},
	{nodeSpec: "z0jndjt42op2/23", canonical: "z0jndjt42op2/23", specType: "node_prev", nodeId: 5486787, nodeIdSize: 23, nodePrev: 0x3ff596748ea77186},
	{nodeSpec: "f2bembkd4zrb/1", canonical: "f2bembkd4zrb/1", specType: "node_prev", nodeId: 1, nodeIdSize: 1, nodePrev: 0x1b844eb5d1aebb07},
	{nodeSpec: "f2bembkd4zrb/8", canonical: "f2bembkd4zrb/8", specType: "node_prev", nodeId: 174, nodeIdSize: 8, nodePrev: 0x1b844eb5d1aebb07},
	{nodeSpec: "f2bembkd4zrb/16", canonical: "f2bembkd4zrb/16", specType: "node_prev", nodeId: 44731, nodeIdSize: 16, nodePrev: 0x1b844eb5d1aebb07},
	{nodeSpec: "f2bembkd4zrb/23", canonical: "f2bembkd4zrb/23", specType: "node_prev", nodeId: 5725571, nodeIdSize: 23, nodePrev: 0x1b844eb5d1aebb07},
	{nodeSpec: "mkg0fd5p76pp/1", canonical: "mkg0fd5p76pp/1", specType: "node_prev", nodeId: 0, nodeIdSize: 1, nodePrev: 0x29391373ab449abd},
	{nodeSpec: "mkg0fd5p76pp/8", canonical: "mkg0fd5p76pp/8", specType: "node_prev", nodeId: 68, nodeIdSize: 8, nodePrev: 0x29391373ab449abd},
	{nodeSpec: "mkg0fd5p76pp/16", canonical: "mkg0fd5p76pp/16", specType: "node_prev", nodeId: 17562, nodeIdSize: 16, nodePrev: 0x29391373ab449abd},
	{nodeSpec: "mkg0fd5p76pp/23", canonical: "mkg0fd5p76pp/23", specType: "node_prev", nodeId: 2248030, nodeIdSize: 23, nodePrev: 0x29391373ab449abd},
}
