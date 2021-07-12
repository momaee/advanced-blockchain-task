package scaffold

import (
	"backend_task/conf"
	grpc_ethereum "backend_task/interface/ethereum"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logrusorgru/aurora"
)

type Scaffold interface {
	server() error

	commissioning() error

	Close()

	SetConfigPath(string)

	NormalStart() error
}

type skeleton struct {
	params   *conf.AppConfig
	listener net.Listener
	context  context.Context
	cancel   context.CancelFunc
}

func Prepare(ctx context.Context, cancelFunc context.CancelFunc) Scaffold {

	s := &skeleton{params: conf.GetAppConfig(), context: ctx, cancel: cancelFunc}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGSTOP)

	go func() {
		for {
			select {
			case <-c:
				s.Close()
			}
		}
	}()

	return s
}

func (s *skeleton) NormalStart() error {

	if err := s.commissioning(); err != nil {
		return err
	}

	if err := s.server(); err != nil {

		return err
	}

	return nil
}

func (s *skeleton) SetConfigPath(path string) {
	conf.GetAppConfig().SetPath(path)
}

func (s *skeleton) eth() error {

	eth, err := grpc_ethereum.NewETHClient(s.context)
	if err != nil {
		return err
	}

	if err := eth.FetchFromTheGraph(); err != nil {
		fmt.Println(aurora.Red(err))
		return err
	}

	for {
		// the program should never exit from run
		if err := eth.ListenToEvents(); err != nil {
			fmt.Println(aurora.Red(err))
		}

		// if program exited from Run, wait for 10 seconds and try again
		time.Sleep(time.Second * 10)
	}
}
