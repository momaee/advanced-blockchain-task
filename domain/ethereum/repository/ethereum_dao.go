package repository

import (
	"backend_task/clients/the_graph"
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type (
	Ethereum interface {
		HandleLog(vLog types.Log) error
		HandleUnprocessedLogs() error
		DbHandleMint(event *the_graph.Event) error
		DbHandleRedeem(event *the_graph.Event) error
	}

	ethereum struct {
		cancel context.CancelFunc
		ctx    context.Context
		mu     sync.Mutex
	}
)

var eth *ethereum

var unprocessedLogs = make(map[string]types.Log)

var (
	MINT   = common.HexToHash("0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f")
	REDEEM = common.HexToHash("0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929")
)

func NewEthereumObject(ctx context.Context) Ethereum {
	if eth == nil {
		eth = new(ethereum)
		eth.ctx, eth.cancel = context.WithCancel(ctx)
	}
	return eth
}
