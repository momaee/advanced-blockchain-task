package grpc_backend

import (
	backend_pb "backend_task/api/pb/commons"
	"backend_task/domain/backend/repository"
	"context"
)

type Server struct {
	backend_pb.UnimplementedBackendServer
}

func (s *Server) GetSupplies(ctx context.Context, in *backend_pb.Request) (*backend_pb.Response, error) {
	return repository.NewBackendObject().GetSupplies(ctx, in)
}
