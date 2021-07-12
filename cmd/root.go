package cmd

import (
	"backend_task/app/scaffold"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CommandLine interface {
	run(cmd *cobra.Command, args []string)
}
type command struct{}

var (
	path    string
	Runner  CommandLine = &command{}
	rootCmd             = &cobra.Command{
		Use:   "backed-task",
		Short: "Mohammad backend task",
		Long: `Mohammad backend task that retrieves events from Compound v2 smart contract, 
		stores them in db and provides rest and grpc api.`,
		Run: Runner.run,
	}
	s scaffold.Scaffold

	systemwideContext context.Context
	cancelFunc        context.CancelFunc
)

func init() {
	systemwideContext, cancelFunc = context.WithCancel(context.Background())
	s = scaffold.Prepare(systemwideContext, cancelFunc)
	rootCmd.Flags().StringVarP(&path, "path", "p", "./conf/", "Set base path for the config file. Default is ./conf ")
	s.SetConfigPath(path)
}

func Execute() error {
	return rootCmd.Execute()
}

func (c *command) run(cmd *cobra.Command, args []string) {
	if err := s.NormalStart(); err != nil {
		fmt.Println("error: ", err)
	}
}
