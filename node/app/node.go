package app

import (
	"errors"

	"source.quilibrium.com/quilibrium/monorepo/node/consensus"
	"source.quilibrium.com/quilibrium/monorepo/node/execution"
	"source.quilibrium.com/quilibrium/monorepo/node/execution/ceremony"
)

type Node struct {
	execEngines map[string]execution.ExecutionEngine
	engine      consensus.ConsensusEngine
}

func newNode(
	ceremonyExecutionEngine *ceremony.CeremonyExecutionEngine,
	engine consensus.ConsensusEngine,
) (*Node, error) {
	if engine == nil {
		return nil, errors.New("engine must not be nil")
	}

	execEngines := make(map[string]execution.ExecutionEngine)
	if ceremonyExecutionEngine != nil {
		execEngines[ceremonyExecutionEngine.GetName()] = ceremonyExecutionEngine
	}

	return &Node{
		execEngines,
		engine,
	}, nil
}

func (n *Node) Start() {
	err := <-n.engine.Start()
	if err != nil {
		panic(err)
	}

	// TODO: add config mapping to engine name/frame registration
	for _, e := range n.execEngines {
		n.engine.RegisterExecutor(e, 0)
	}
}

func (n *Node) Stop() {
	err := <-n.engine.Stop(false)
	if err != nil {
		panic(err)
	}
}