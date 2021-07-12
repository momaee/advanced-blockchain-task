package the_graph

import (
	"context"
)

type (
	TheGraph interface {
		GetMintEvent(ctx context.Context, hash, logIndex string) (*Event, error)
		GetRedeemEvent(ctx context.Context, hash, logIndex string) (*Event, error)
		GetLastDayMintEvents(ctx context.Context, contractAddress string) ([]*Event, error)
		GetLastDayRedeemEvents(ctx context.Context, contractAddress string) ([]*Event, error)
		GetAllMarkets(ctx context.Context) ([]*Market, error)
	}

	theGraph struct {
	}
)

type Payload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type Event struct {
	Id               string `json:"id"`
	Amount           string `json:"amount"`
	To               string `json:"to"`
	From             string `json:"from"`
	BlockNumber      int    `json:"blockNumber"`
	BlockTime        int    `json:"blockTime"`
	UnderlyingAmount string `json:"underlyingAmount"`
}

type MintEventResponse struct {
	Data *MintEventData `json:"data"`
}

type MintEventData struct {
	MintEvents []*Event `json:"mintEvents"`
}

type RedeemEventResponse struct {
	Data *RedeemEventData `json:"data"`
}

type RedeemEventData struct {
	RedeemEvents []*Event `json:"redeemEvents"`
}

type MarketResponse struct {
	Data *MarketData `json:"data"`
}

type MarketData struct {
	Markets []*Market `json:"markets"`
}

type Market struct {
	Id string `json:"id"`
}

func New() TheGraph {
	return new(theGraph)
}
