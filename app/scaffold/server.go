package scaffold

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	backend "backend_task/api/pb/commons"
	grpc_backend "backend_task/interface/backend"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/logrusorgru/aurora"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *skeleton) server() error {
	var err error

	s.listener, err = net.Listen("tcp",
		net.JoinHostPort(s.params.Server.Grpc.Host, s.params.Server.Grpc.Port))

	if err != nil {

		return fmt.Errorf("cannot open grpc %v\n", aurora.Red(err))
	}

	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			// we can add middlewares here
			// grpc_auth.StreamServerInterceptor(middleware.AuthenticateToken),
		),
	),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				// grpc_auth.UnaryServerInterceptor(middleware.AuthenticateToken),
			),
		),
	)
	go func() {
		reflection.Register(grpcServer)
		grpcServer.Serve(s.listener)

	}()
	// register your services here
	backend.RegisterBackendServer(grpcServer, &grpc_backend.Server{})

	conn, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort(s.params.Server.Grpc.Host, s.params.Server.Grpc.Port),
		grpc.WithInsecure(),
	)

	if err != nil {
		return err
	}

	router := runtime.NewServeMux()

	if err = backend.RegisterBackendHandler(context.Background(), router, conn); err != nil {
		log.Fatalf("Failed to connect to register gateway: %v\n", err)
	}

	return http.ListenAndServe(net.JoinHostPort(s.params.Server.Rest.Host, s.params.Server.Rest.Port),
		wsproxy.WebsocketProxy(router))

}
