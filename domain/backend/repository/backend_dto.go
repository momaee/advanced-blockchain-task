package repository

import (
	backend_pb "backend_task/api/pb/commons"
	"backend_task/clients/redis"
	"backend_task/domain/ethereum/repository"
	"context"
	"fmt"
	"time"

	go_redis "github.com/go-redis/redis/v8"
)

func (b *backend) GetSupplies(ctx context.Context, in *backend_pb.Request) (*backend_pb.Response, error) {
	var value float64

	db, err := redis.Storage.GetDB()
	if err != nil {
		return nil, err
	}

	supplies := []*backend_pb.HourlySupply{}
	hour := time.Now().Unix() - time.Now().Unix()%(60*60)

	for i := hour - 23*60*60; i <= hour; i += 60 * 60 {
		key := repository.MakeRedisKey(in.ContractAddress, int(i))
		err = db.Get(ctx, key, &value)
		//todo if everything is ok this should not happen
		if err == go_redis.Nil {
			value = 0
			supplies = append(supplies, &backend_pb.HourlySupply{
				Timestamp:   i,
				TotalSupply: fmt.Sprintf("%v", value),
			})
			continue
		}
		if err != nil {
			return nil, err
		}
		supplies = append(supplies, &backend_pb.HourlySupply{
			Timestamp:   i,
			TotalSupply: fmt.Sprintf("%v", value),
		})
	}

	return &backend_pb.Response{
		Supplies: supplies,
	}, nil
}
