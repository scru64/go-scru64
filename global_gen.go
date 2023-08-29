package scru64

import (
	"fmt"
	"os"
	"sync"
)

// The gateway object that forwards supported method calls to the process-wide
// global generator.
//
// The global generator reads the node configuration from the `SCRU64_NODE_SPEC`
// environment variable by default, and it panics if it fails to read a
// well-formed node spec string (e.g., `"42/8"`, `"0xb00/12"`,
// `"0u2r85hm2pt3/16"`) when a generator method is first called. See also
// [NodeSpec] for the node spec string format.
var GlobalGenerator interface {
	// Initializes the global generator, if not initialized, with the node spec
	// passed.
	//
	// This method configures the global generator with the argument only when the
	// global generator is not yet initialized. Otherwise, it preserves the
	// existing configuration.
	//
	// This method return `true` if this method configures the global generator or
	// `false` if it preserves the existing configuration.
	Initialize(nodeSpec NodeSpec) bool

	// Calls `Generator.Generate` of the global generator.
	Generate() (Id, error)

	// Calls `Generator.GenerateOrSleep` of the global generator.
	GenerateOrSleep() Id

	// Calls `Generator.NodeId` of the global generator.
	NodeId() uint32

	// Calls `Generator.NodeIdSize` of the global generator.
	NodeIdSize() uint8

	// Calls `Generator.NodeSpec` of the global generator.
	NodeSpec() NodeSpec
} = &globalGeneratorInner{}

// The lazy initialization holder type of the global generator.
type globalGeneratorInner struct {
	once  sync.Once
	inner *Generator
}

func (g *globalGeneratorInner) get() *Generator {
	g.once.Do(func() {
		nodeSpec, err := ParseNodeSpec(os.Getenv("SCRU64_NODE_SPEC"))
		if err != nil {
			panic(fmt.Errorf(
				"scru64: could not read config from SCRU64_NODE_SPEC env var: %w", err))
		}
		g.inner = NewGenerator(nodeSpec)
	})
	return g.inner
}

func (g *globalGeneratorInner) Initialize(nodeSpec NodeSpec) bool {
	initialized := false
	g.once.Do(func() {
		g.inner = NewGenerator(nodeSpec)
		initialized = true
	})
	return initialized
}

func (g *globalGeneratorInner) Generate() (Id, error) {
	return g.get().Generate()
}

func (g *globalGeneratorInner) GenerateOrSleep() Id {
	return g.get().GenerateOrSleep()
}

func (g *globalGeneratorInner) NodeId() uint32 {
	return g.get().NodeId()
}

func (g *globalGeneratorInner) NodeIdSize() uint8 {
	return g.get().NodeIdSize()
}

func (g *globalGeneratorInner) NodeSpec() NodeSpec {
	return g.get().NodeSpec()
}
