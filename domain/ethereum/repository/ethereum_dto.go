package repository

import (
	"backend_task/clients/redis"
	"backend_task/clients/the_graph"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/logrusorgru/aurora"
)

func (e *ethereum) HandleUnprocessedLogs() error {
	for _, log := range unprocessedLogs {
		if err := e.HandleLog(log); err != nil {
			return err
		}
	}
	return nil
}

func (e *ethereum) HandleLog(vLog types.Log) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// todo handle this better
	if len(vLog.Topics) != 1 {
		fmt.Println(aurora.Red("This case never should happen. log:"), vLog)
		return nil
	}

	topic := vLog.Topics[0]

	if topic == MINT {
		if err := handleMint(e.ctx, vLog); err != nil {
			return err
		}
	}
	if topic == REDEEM {
		if err := handleRedeem(e.ctx, vLog); err != nil {
			return err
		}
	}

	return nil
}

func handleMint(ctx context.Context, vLog types.Log) error {
	// I check my event with the graph. this has a overhead but I get the timestamp and float amount
	// Actually I'm not very ok with this work but I did it.
	event, err := the_graph.New().GetMintEvent(ctx, vLog.TxHash.Hex(), strconv.Itoa(int(vLog.Index)))
	if err != nil {
		return err
	}

	if event == nil {
		addUnprocessedLog(vLog)
		return nil
	}

	if !vLog.Removed {
		if err := dbHandleMint(ctx, vLog, event); err != nil {
			return err
		}
	} else {
		if err := dbUndoMint(ctx, vLog, event); err != nil {
			return err
		}
	}

	checkUnprocessedLogs(vLog)
	return nil
}

func handleRedeem(ctx context.Context, vLog types.Log) error {
	// I check my event with the graph. this has a overhead but I get the timestamp and float amount
	// Actually I'm not very ok with this work but I did it.
	event, err := the_graph.New().GetRedeemEvent(ctx, vLog.TxHash.Hex(), strconv.Itoa(int(vLog.Index)))
	if err != nil {
		return err
	}

	if event == nil {
		addUnprocessedLog(vLog)
		return nil
	}

	if !vLog.Removed {
		if err := dbHandleRedeem(ctx, vLog, event); err != nil {
			return err
		}
	} else {
		if err := dbUndoRedeem(ctx, vLog, event); err != nil {
			return err
		}
	}

	checkUnprocessedLogs(vLog)
	return nil
}

// all processing finished, if the log exists in unprocessedlogs, delete it.
func checkUnprocessedLogs(log types.Log) {
	delete(unprocessedLogs, createUnprocessedLogKey(log))
}

func createUnprocessedLogKey(vLog types.Log) string {
	return strings.ToLower(vLog.TxHash.Hex()) + "-" + strconv.Itoa(int(vLog.Index)) + "-" + strconv.FormatBool(vLog.Removed)
}

func addUnprocessedLog(vLog types.Log) {
	_, ok := unprocessedLogs[createUnprocessedLogKey(vLog)]
	if !ok {
		unprocessedLogs[createUnprocessedLogKey(vLog)] = vLog
	}
}

func dbHandleMint(ctx context.Context, log types.Log, event *the_graph.Event) error {
	hour := event.BlockTime - event.BlockTime%(60*60)
	key := MakeRedisKey(log.Address.Hex(), hour)

	fAmount, err := strconv.ParseFloat(event.UnderlyingAmount, 64)
	if err != nil {
		return err
	}

	if err := supply(ctx, key, fAmount); err != nil {
		return err
	}

	// fmt.Println("hash", vLog.TxHash.Hex())
	// fmt.Println("data minter", common.BytesToAddress(vLog.Data[0:32]))
	// fmt.Println("data mintAmount", big.NewInt(0).SetBytes(vLog.Data[32:64])) //supply
	// fmt.Println("data mintTokens", big.NewInt(0).SetBytes(vLog.Data[64:96]))
	return nil
}

func (e *ethereum) DbHandleMint(event *the_graph.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	hour := event.BlockTime - event.BlockTime%(60*60)
	key := MakeRedisKey(event.From, hour)
	fAmount, err := strconv.ParseFloat(event.UnderlyingAmount, 64)
	if err != nil {
		return err
	}
	if err := supply(e.ctx, key, fAmount); err != nil {
		return err
	}
	return nil
}

func dbHandleRedeem(ctx context.Context, log types.Log, event *the_graph.Event) error {
	hour := event.BlockTime - event.BlockTime%(60*60)
	key := MakeRedisKey(log.Address.Hex(), hour)

	fAmount, err := strconv.ParseFloat(event.UnderlyingAmount, 64)
	if err != nil {
		return err
	}

	if err := withdraw(ctx, key, fAmount); err != nil {
		return err
	}
	// event.UnderlyingAmount withdraw
	// fmt.Println("hash", vLog.TxHash.Hex())
	// fmt.Println("data redeemer", common.BytesToAddress(vLog.Data[0:32]))
	// fmt.Println("data redeemAmount", big.NewInt(0).SetBytes(vLog.Data[32:64])) //withdraw
	// fmt.Println("data redeemTokens", big.NewInt(0).SetBytes(vLog.Data[64:96]))
	return nil
}
func (e *ethereum) DbHandleRedeem(event *the_graph.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	hour := event.BlockTime - event.BlockTime%(60*60)
	key := MakeRedisKey(event.From, hour)
	fAmount, err := strconv.ParseFloat(event.UnderlyingAmount, 64)
	if err != nil {
		return err
	}
	if err := withdraw(e.ctx, key, fAmount); err != nil {
		return err
	}
	return nil
}

func dbUndoMint(ctx context.Context, log types.Log, event *the_graph.Event) error {
	return dbHandleRedeem(ctx, log, event)
}

func dbUndoRedeem(ctx context.Context, log types.Log, event *the_graph.Event) error {
	return dbHandleMint(ctx, log, event)
}

func supply(ctx context.Context, key string, amount float64) error {
	db, err := redis.Storage.GetDB()
	if err != nil {
		return err
	}

	var value float64
	err = db.Get(ctx, key, &value)
	if err == go_redis.Nil {
		value = amount
		if err := db.Set(ctx, key, value, 0); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}

	value += amount
	if err := db.Set(ctx, key, value, 0); err != nil {
		return err
	}

	return nil
}

func withdraw(ctx context.Context, key string, amount float64) error {
	db, err := redis.Storage.GetDB()
	if err != nil {
		return err
	}
	var value float64
	err = db.Get(ctx, key, &value)
	if err == go_redis.Nil {
		// not exist create it
		value = -amount
		if err := db.Set(ctx, key, value, 0); err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	value -= amount
	if err := db.Set(ctx, key, value, 0); err != nil {
		return err
	}
	return nil
}

func MakeRedisKey(address string, time int) string {

	return strings.ToLower(address) + "-" + strconv.Itoa(time)
}
