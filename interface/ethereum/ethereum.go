package grpc_ethereum

import (
	"context"
	"fmt"
	"html"
	"math/big"
	"time"

	"backend_task/clients/the_graph"
	"backend_task/domain/ethereum/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/logrusorgru/aurora"
)

type (
	Eth interface {
		ListenToEvents() error
		FetchFromTheGraph() error
	}
	eth struct {
		*ethclient.Client
		cancel  context.CancelFunc
		ctx     context.Context
		repEth  repository.Ethereum
		markets []common.Address
	}
)

var client *eth

func NewETHClient(ctx context.Context) (Eth, error) {
	if client == nil {
		client = new(eth)
		client.ctx, client.cancel = context.WithCancel(ctx)
		client.repEth = repository.NewEthereumObject(ctx)
		if err := client.fetchMarkets(); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (e *eth) fetchMarkets() error {
	markets, err := the_graph.New().GetAllMarkets(e.ctx)
	if err != nil {
		return nil
	}
	for _, market := range markets {
		e.markets = append(e.markets, common.HexToAddress(market.Id))
	}
	return nil
}

func (e *eth) connect() error {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/8d810610fe7741cc9753cbaafb1f000c")
	if err != nil {
		return fmt.Errorf("cannot connect to the client %v", err)
	}
	e.Client = client
	fmt.Printf("%s connected to geth client.\n", aurora.Green(html.UnescapeString("&#x2705;")))
	return nil
}

// func (e *eth) FetchFromNetwork() error {
// 	if err := e.connect(); err != nil {
// 		return err
// 	}

// 	query := ethereum.FilterQuery{
// 		FromBlock: big.NewInt(12800000),
// 		Topics: [][]common.Hash{
// 			{repository.MINT, repository.REDEEM},
// 		},
// 		Addresses: []common.Address{repository.CETH},
// 	}

// 	logs, err := e.FilterLogs(context.Background(), query)
// 	if err != nil {
// 		return err
// 	}

// 	for _, vLog := range logs {
// 		fmt.Println("log", vLog)
// 	}
// 	return nil
// }

func (e *eth) ListenToEvents() error {
	if err := e.connect(); err != nil {
		return err
	}

	query := ethereum.FilterQuery{
		// just get events that topic[0]=MINT or topic[o]=REDEEM
		Topics: [][]common.Hash{
			{repository.MINT, repository.REDEEM},
		},
		Addresses: e.markets,
	}

	logs := make(chan types.Log)

	sub, err := e.SubscribeFilterLogs(e.ctx, query, logs)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {

		case <-e.ctx.Done():
			return fmt.Errorf("error in proccessing events")

		case err := <-sub.Err():
			return err

		case <-ticker.C:
			if err := e.repEth.HandleUnprocessedLogs(); err != nil {
				return err
			}

		case vLog := <-logs:
			printEvent(vLog)
			if err := e.repEth.HandleLog(vLog); err != nil {
				return err
			}
		}
	}
}

//todo handle rewriting
//todo handle lastblock of the graph
//todo this process is too slow. solution: 1-calculate groups then query to db 2-fetch data from graph in one query
func (e *eth) FetchFromTheGraph() error {
	for _, market := range e.markets {
		mintEvents, err := the_graph.New().GetLastDayMintEvents(e.ctx, market.Hex())
		if err != nil {
			return err
		}
		//todo this loop is stupid, but redis is fast enough
		for _, event := range mintEvents {
			e.repEth.DbHandleMint(event)
		}

		redeemEvents, err := the_graph.New().GetLastDayRedeemEvents(e.ctx, market.Hex())
		if err != nil {
			return err
		}
		//todo this loop is stupid, but redis is fast enough
		for _, event := range redeemEvents {
			e.repEth.DbHandleRedeem(event)
		}
		fmt.Println(aurora.Green("last 24 hours fetched from the graph."), aurora.Green("market:"), aurora.Green(market))
	}
	return nil
}

func printEvent(vLog types.Log) {
	if len(vLog.Topics) != 1 {
		fmt.Println(aurora.Red("This case never should happen. log:"), vLog)
		return
	}

	topic := vLog.Topics[0]

	if topic == repository.MINT {
		fmt.Println(aurora.Green("new mint event"))
		fmt.Println("contract address", aurora.Green(vLog.Address.Hex()))
		fmt.Println("block number", aurora.Green(vLog.BlockNumber))
		fmt.Println("block hash", aurora.Green(vLog.BlockHash.Hex()))
		fmt.Println("tx hash", aurora.Green(vLog.TxHash.Hex()))
		fmt.Println("minter", aurora.Green(common.BytesToAddress(vLog.Data[0:32])))
		fmt.Println("mintAmount", aurora.Green(big.NewInt(0).SetBytes(vLog.Data[32:64])))
		fmt.Println("mintTokens", aurora.Green(big.NewInt(0).SetBytes(vLog.Data[64:96])))
		fmt.Println("")
	}
	if topic == repository.REDEEM {
		fmt.Println(aurora.Blue("new redeem event"))
		fmt.Println("contract address", aurora.Blue(vLog.Address.Hex()))
		fmt.Println("block number", aurora.Blue(vLog.BlockNumber))
		fmt.Println("block hash", aurora.Blue(vLog.BlockHash.Hex()))
		fmt.Println("tx hash", aurora.Blue(vLog.TxHash.Hex()))
		fmt.Println("redeemer", aurora.Blue(common.BytesToAddress(vLog.Data[0:32])))
		fmt.Println("redeemAmount", aurora.Blue(big.NewInt(0).SetBytes(vLog.Data[32:64])))
		fmt.Println("redeemTokens", aurora.Blue(big.NewInt(0).SetBytes(vLog.Data[64:96])))
		fmt.Println("")
	}
}
