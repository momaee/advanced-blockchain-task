package repository

import (
	backend_pb "backend_task/api/pb/commons"
	"context"
)

type (
	Backend interface {
		GetSupplies(ctx context.Context, in *backend_pb.Request) (*backend_pb.Response, error)
	}

	backend struct {
	}
)

func NewBackendObject() Backend {
	return new(backend)
}
